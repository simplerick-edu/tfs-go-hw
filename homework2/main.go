package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

type Op struct {
	OpSign       interface{} `json:"type,omitempty"`
	Value        interface{} `json:"value,omitempty"`
	ID           interface{} `json:"id,omitempty"`
	CreationTime interface{} `json:"created_at,omitempty"`
}

type Transaction struct {
	Company   string `json:"company,omitempty"`
	Operation Op     `json:"operation,omitempty"`
	Op
}

func (r *Transaction) SearchInOperation() {
	op := r.Operation
	if op.ID != nil {
		r.ID = op.ID
	}
	if op.OpSign != nil {
		r.OpSign = op.OpSign
	}
	if op.Value != nil {
		r.Value = op.Value
	}
	if op.CreationTime != nil {
		r.CreationTime = op.CreationTime
	}
}

func (r *Transaction) Before(r2 Transaction) bool {
	time1, _ := convertToTime(r.CreationTime)
	time2, _ := convertToTime(r2.CreationTime)
	return time1.Before(time2)
}

// check id is valid
func checkID(obj interface{}) (interface{}, bool) {
	switch v := obj.(type) {
	case string:
		return v, true
	case float64:
		return floatToInt(v)
	}
	return obj, false
}

// determines sign {+1, -1} of operation type
func convertToSign(obj interface{}) (int, bool) {
	if s, ok := obj.(string); ok {
		switch s {
		case "income":
			return 1, true
		case "+":
			return 1, true
		case "outcome":
			return -1, true
		case "-":
			return -1, true
		}
	}
	return 0, false
}

func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}

// converts float to int, ok if conversion lossless
func floatToInt(x float64) (int, bool) {
	y := int(x)
	if float64(y) == x {
		return y, true
	}
	return y, false
}

// converts strings and floats to non-negative int
func convertToInt(obj interface{}) (int, bool) {
	var defaultValue int
	// string
	if v, ok := obj.(string); ok {
		converted, err := strconv.Atoi(v)
		return abs(converted), (err == nil)
	}
	// float
	if v, ok := obj.(float64); ok {
		if x, ok := floatToInt(v); ok {
			return abs(x), true
		}
	}
	return defaultValue, false
}

// convert to time format
func convertToTime(obj interface{}) (time.Time, bool) {
	RFC3339 := "2006-01-02T15:04:05Z07:00"
	var t time.Time
	if s, ok := obj.(string); ok {
		t, err := time.Parse(RFC3339, s)
		return t, err == nil
	}
	return t, false
}

type ConsolidatedReport struct {
	Company               string        `json:"company"`
	ValidOpOperationCount int           `json:"valid_operations_count"`
	Balance               int           `json:"balance"`
	InvalidOperations     []interface{} `json:"invalid_operations,omitempty"`
}

func New(companyName string) ConsolidatedReport {
	return ConsolidatedReport{
		Company:           companyName,
		InvalidOperations: make([]interface{}, 0),
	}
}

func (r *ConsolidatedReport) addTransaction(rec Transaction) error {
	if r.Company == "" {
		return errors.New("Company name is empty")
	}
	if r.Company != rec.Company {
		return errors.New("Company names do not match")
	}
	if id, ok := checkID(rec.ID); ok {
		if _, ok := convertToTime(rec.CreationTime); ok {
			opSign, ok1 := convertToSign(rec.OpSign)
			value, ok2 := convertToInt(rec.Value)
			if ok1 && ok2 {
				r.Balance += opSign * value
				r.ValidOpOperationCount++
			} else {
				r.InvalidOperations = append(r.InvalidOperations, id)
			}
		}
	}
	return errors.New("Transaction is invalid")
}

func readData() ([]byte, error) {
	var filePathFlag = flag.String("file", "", "path to file")
	flag.Parse()
	filePathEnv, okEnv := os.LookupEnv("FILE")
	switch {
	case *filePathFlag != "":
		data, err := os.ReadFile(*filePathFlag)
		return data, err
	case okEnv:
		data, err := os.ReadFile(filePathEnv)
		return data, err
	default:
		data, err := io.ReadAll(os.Stdin)
		return data, err
	}
}

func writeData(path string, obj interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	if err = enc.Encode(obj); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func main() {
	data, readErr := readData()
	if readErr != nil {
		fmt.Println(fmt.Errorf("Error when trying to read a file: %w", readErr))
	}
	var transactions []Transaction
	err := json.Unmarshal(data, &transactions)
	if err != nil {
		fmt.Println(fmt.Errorf("Error when unmarshaling: %w", err))
	}

	for i := range transactions {
		transactions[i].SearchInOperation()
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Before(transactions[j])
	})

	var results = []ConsolidatedReport{}
	var resultsIdx = map[string]int{}

	for _, transaction := range transactions {
		companyName := transaction.Company
		if _, ok := resultsIdx[companyName]; !ok {
			resultsIdx[companyName] = len(results)
			results = append(results, New(companyName))
		}
		_ = results[resultsIdx[companyName]].addTransaction(transaction)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Company < results[j].Company
	})
	writeErr := writeData("out.json", results)
	if writeErr != nil {
		fmt.Println(fmt.Errorf("Error when trying to write to a file: %w", writeErr))
	}
}
