package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type httpResponse struct {
	IDlist []string `json:"id_list"`
	Passed bool     `json:"passed"`
}

func main() {
	myScrapeInfo := ScrapeInfo{
		baseURL:       "https://daigakujc.jp",
		preScrapePath: "/pal.php?u=31&h=24",
		examNumber:    "0692",
		examCategory:  "工学部",
	}
	passedIDs, iHasPassed := Scrape(myScrapeInfo)
	res := httpResponse{
		IDlist: passedIDs,
		Passed: iHasPassed,
	}

	jsonHandler := func(rw http.ResponseWriter, req *http.Request) {
		responseBody, err := json.Marshal(res)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("content-Type", "application/json")
		fmt.Fprint(rw, string(responseBody))
	}

	http.HandleFunc("/", jsonHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, nil)
}
