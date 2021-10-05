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

type Transaction struct {
	Company      string
	OpSign       int
	Value        int
	ID           interface{}
	CreationTime time.Time
	Success      bool
	ValidOp      bool
}

func (r *Transaction) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if op, ok := raw["operation"].(map[string]interface{}); ok {
		for k, v := range op {
			raw[k] = v
		}
	}
	// assigning variables
	if r.Company, r.Success = raw["company"].(string); r.Success {
		if r.ID, r.Success = checkID(raw["id"]); r.Success {
			r.CreationTime, r.Success = convertToTime(raw["created_at"])
		}
	}
	if r.OpSign, r.ValidOp = convertToSign(raw["type"]); r.ValidOp {
		r.Value, r.ValidOp = convertToInt(raw["value"])
	}
	return nil
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
	if r.Company != rec.Company {
		return errors.New("Company names do not match")
	}
	if rec.ValidOp {
		r.Balance += rec.OpSign * rec.Value
		r.ValidOpOperationCount++
	} else {
		r.InvalidOperations = append(r.InvalidOperations, rec.ID)
	}
	return nil
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
	_ = json.Unmarshal(data, &transactions)

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreationTime.Before(transactions[j].CreationTime)
	})

	var results = []ConsolidatedReport{}
	var resultsIdx = map[string]int{}

	for _, transaction := range transactions {
		if transaction.Success {
			companyName := transaction.Company
			if _, ok := resultsIdx[companyName]; !ok {
				resultsIdx[companyName] = len(results)
				results = append(results, New(companyName))
			}
			_ = results[resultsIdx[companyName]].addTransaction(transaction)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Company < results[j].Company
	})
	writeErr := writeData("out.json", results)
	if writeErr != nil {
		fmt.Println(fmt.Errorf("Error when trying to write to a file: %w", writeErr))
	}
}
