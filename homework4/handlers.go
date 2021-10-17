package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"io"
	"net/http"
	"time"
)

var (
	signKey = []byte("secret_key")
)

const (
	cookieAuth = "auth"
	userID     = "ID"
)

// map with usernames and passwords
var users = map[string]string{}

type Message struct {
	Username string `json:"user"`
	Text     string `json:"text"`
}

type Messages []Message

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserPair struct {
	user1, user2 string
}

// storage for messages
var chatMessages = Messages{}
var userMessages = map[UserPair]*Messages{}

func createToken(user string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))
	t.Claims = jwt.StandardClaims{
		Id:        user,
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
	}
	return t.SignedString(signKey)
}

// authentification
func Login(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var creds Credentials
	err = json.Unmarshal(d, &creds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// generate new token
	token, err := createToken(creds.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:  cookieAuth,
		Value: token,
	})

}

// registration
func Signup(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var creds Credentials
	err = json.Unmarshal(d, &creds)

	if err != nil || creds.Username == "" || creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, ok := users[creds.Username]; !ok {
		users[creds.Username] = creds.Password
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

// authorization
func Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieAuth)

		switch err {
		case nil:
		case http.ErrNoCookie:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tknStr := c.Value
		claims := &jwt.StandardClaims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return signKey, nil
		})

		switch err {
		case nil:
		case jwt.ErrSignatureInvalid:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		idCtx := context.WithValue(r.Context(), userID, claims.Id)
		handler.ServeHTTP(w, r.WithContext(idCtx))
	}
	return http.HandlerFunc(fn)
}

func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	getMessages(chatMessages, w, r)
}

func PostChatMessages(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(userID).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postMessages(&chatMessages, id, w, r)
}

func GetUserMessages(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(userID).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id2 := chi.URLParam(r, "id")
	m := userMessages[UserPair{id, id2}]
	if m == nil {
		m = &Messages{}
		userMessages[UserPair{id, id2}] = m
		userMessages[UserPair{id2, id}] = m
	}
	getMessages(*m, w, r)
}

func PostUserMessages(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(userID).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id2 := chi.URLParam(r, "id")
	m := userMessages[UserPair{id, id2}]
	if m == nil {
		m = &Messages{}
		userMessages[UserPair{id, id2}] = m
		userMessages[UserPair{id2, id}] = m
	}
	postMessages(m, id, w, r)
}

func getMessages(m Messages, w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(userID).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func postMessages(m *Messages, userId string, w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var text string
	err = json.Unmarshal(d, &text)
	fmt.Println(err)
	if err != nil || text == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	*m = append(*m, Message{Username: userId, Text: text})
}
