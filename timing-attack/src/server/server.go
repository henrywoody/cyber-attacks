package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/data/vulnerable", handleRequestVulnerable)
	http.HandleFunc("/data/secure", handleRequestSecure)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

const password string = "password"
const data string = "Secret information"

type passwordPayload struct {
	Password string `json:"password"`
}

func handleRequestVulnerable(w http.ResponseWriter, r *http.Request) {
	var payload passwordPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	success := vulnerableComparePassword(payload.Password)
	if !success {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func vulnerableComparePassword(givenPassword string) bool {
	minLen := int(math.Min(float64(len(givenPassword)), float64(len(password))))
	for i := 0; i < minLen; i++ {
		time.Sleep(10 * time.Microsecond)
		if givenPassword[i] != password[i] {
			return false
		}
	}

	if len(givenPassword) != len(password) {
		return false
	}

	return true
}

func handleRequestSecure(w http.ResponseWriter, r *http.Request) {
	var payload passwordPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	success := secureComparePassword(payload.Password)
	if !success {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func secureComparePassword(givenPassword string) bool {
	passwordMatches := true

	minLen := int(math.Min(float64(len(givenPassword)), float64(len(password))))
	for i := 0; i < minLen; i++ {
		time.Sleep(10 * time.Microsecond)
		if givenPassword[i] != password[i] {
			passwordMatches = false
		}
	}

	if !passwordMatches {
		return false
	}

	if len(givenPassword) != len(password) {
		return false
	}

	return true
}
