package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var mostCommonSize = 3
var serverMostCommonWords map[string]int
var serverWords map[string]int

func add(w http.ResponseWriter, req *http.Request) {
	words := strings.Split(req.URL.Query().Get("data"), ",")
	w.WriteHeader(http.StatusOK)
	if len(words) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		serverWords[word]++
		if _, ok := serverMostCommonWords[word]; ok {
			serverMostCommonWords[word]++
		} else {
			if len(serverMostCommonWords) < mostCommonSize {
				serverMostCommonWords[word] = serverWords[word]
			} else {
				smallestWordSize := serverWords[word]
				smallestWord := word
				for commonWord, _ := range serverMostCommonWords {
					if serverMostCommonWords[commonWord] < smallestWordSize {
						smallestWordSize = serverMostCommonWords[commonWord]
						smallestWord = commonWord
					}
				}
				if serverWords[word] > smallestWordSize {
					delete(serverMostCommonWords, smallestWord)
					serverMostCommonWords[word] = serverWords[word]
				}
			}
		}
	}

	_, _ = fmt.Fprintf(w, "Added words %v to the map", words)
}

func status(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(serverMostCommonWords)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(marshal)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func init() {
	fmt.Println("Starting the application...")
	serverWords = make(map[string]int)
	serverMostCommonWords = make(map[string]int, mostCommonSize)
}

func main() {
	http.HandleFunc("/add", add)
	http.HandleFunc("/status", status)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		panic(err)
	}
}
