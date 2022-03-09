package resultInformer

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type ScrapeInfo struct {
	baseURL       string
	preScrapePath string
	examCategory  string
	examNumber    string
}

func Scrape(si ScrapeInfo) ([]string, bool) {
	var (
		baseURL        = si.baseURL
		preScrapePath  = si.preScrapePath
		myExamCategory = si.examCategory
		myExamNumber   = si.examNumber
	)
	var (
		preScrapeURL  = baseURL + preScrapePath
		myCategoryURL string
	)

	href := findHrefOf(myExamCategory, preScrapeURL)

	if href == "not found" {
		return nil, false
	}

	if href[0:5] == "https" {
		myCategoryURL = href
	} else {
		myCategoryURL = baseURL + href
	}

	passedIDs := findPassedIDsFrom(myCategoryURL)
	iHasPassed := false
	for _, id := range passedIDs {
		if id[0:4] == myExamNumber {
			iHasPassed = true
			break
		}
	}

	return passedIDs, iHasPassed
}

func findHrefOf(targetWord, targetURL string) string {
	doc, err := goquery.NewDocument(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	var (
		targetHref string
		isUTF8site bool = false //default
		targetCode string
	)

	doc.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		charset, exists := s.Attr("charset")
		if exists {
			isUTF8site = (charset == "utf-8" || charset == "UTF-8")
		}
		return !exists
		// break if !exists is false, in other words, exists is true.
	})

	if isUTF8site {
		targetCode = targetWord
	} else {
		targetCode = convertUTF8toSjis(targetWord)
	}

	foundTarget := false
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		isTargetTag := (s.Text() == targetCode)

		if exists && isTargetTag {
			targetHref = href
			foundTarget = true
		}
		return !isTargetTag
		// break if !isTagetTag is false, in other words, isTargetTag is true.
	})

	if !foundTarget {
		targetHref = "not found"
	}
	return targetHref
}

func findPassedIDsFrom(targetURL string) []string {
	doc, err := goquery.NewDocument(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	var (
		wordList  = strings.Split(doc.Text(), "\n")
		passedIDs []string
		idPattern = regexp.MustCompile(`[0-2]\d{3}[A-F]`)
		// (0001〜2569) + (A〜F)
	)
	for _, word := range wordList {
		id := strings.TrimSpace(word)
		if idPattern.MatchString(id) {
			passedIDs = append(passedIDs, id)
		}
	}

	return passedIDs
}

func convertUTF8toSjis(utf8Str string) string {
	encoder := japanese.ShiftJIS.NewEncoder()
	sjisStr, _, err := transform.String(encoder, utf8Str)
	if err != nil {
		log.Fatal(err)
	}
	return sjisStr
}
