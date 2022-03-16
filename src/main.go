package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type httpResponse struct {
	IDlist []string `json:"id_list"`
	Passed bool     `json:"passed"`
}

var myScrapeInfo = ScrapeInfo{
	baseURL:       "https://daigakujc.jp",
	preScrapePath: "/pal.php?u=31&h=24",
	/*
	examCategory:  os.Getenv("MY_EXAM_CATEGORY"),
	examNumber:    os.Getenv("MY_EXAM_NUMBER"),
	*/
	// test
	examCategory: "工学部",
	examNumber: "0692",
}

var execPoint time.Time
var nowFunc func() time.Time

func init() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

	nowFunc = func() time.Time {
		return time.Now().In(jst)
	}
	execPoint = time.Date(2022, 3, 10, 12, 0, 0, 0, jst)
}

/*
func main() {
		for {
			now := nowFunc()
			diffMin := execPoint.Sub(now).Minutes()
			if diffMin < 5 {
				break
			}
			time.Sleep(10 * time.Second)
		}
		for {
			now := nowFunc()
			diffSec := execPoint.Sub(now).Seconds()
			if diffSec < 10 {
				break
			}
			time.Sleep(5 * time.Second)
		}
		for {
			now := nowFunc()
			diffMilisec := execPoint.Sub(now).Milliseconds()
			if diffMilisec < 500 {
				break
			}
			if diffMilisec < 600 {
				time.Sleep(10 * time.Millisecond)
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
		time.Sleep(500 * time.Millisecond)

		var response = getRes()
		for {
			if len(response.IDlist) > 0 {
				break
			}
			time.Sleep(50 * time.Millisecond)
			response = getRes()
		}

		startServer(response)
}
*/

func main() {
	response := getRes()
	//startServer(response)
	fmt.Println(response.IDlist)
}

func startServer(res httpResponse) {
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

func getRes() httpResponse {
	passedIDs, iHasPassed := scrape(myScrapeInfo)
	res := httpResponse{
		IDlist: passedIDs,
		Passed: iHasPassed,
	}
	return res
}
