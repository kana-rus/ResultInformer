package main

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

func scrape(si ScrapeInfo) ([]string, bool) {
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
		isUTF8site bool = false //default
		targetCode string
	)
	doc.Find("meta").EachWithBreak(func(i int, s *goquery.Selection) bool {
		charset, exists := s.Attr("charset")
		if exists {
			isUTF8site = (charset == "utf-8")
		}
		return !exists
	})
	if isUTF8site {
		targetCode = targetWord
	} else {
		targetCode = convertUTF8toSjis(targetWord)
	}

	var targetHref string
	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		isTargetTag := (s.Text() == targetCode)

		if exists && isTargetTag {
			targetHref = href
		}
		return !isTargetTag
	})

	return targetHref
}

func findPassedIDsFrom(targetURL string) []string {
	doc, err := goquery.NewDocument(targetURL)
	if err != nil {
		log.Fatal(err)
	}
	var (
		idPattern  = regexp.MustCompile(`[0-2]\d{3}[A-F]`)
		targetText string
	)
	doc.Find("font").EachWithBreak(func(i int, s *goquery.Selection) bool {
		str := s.Text()
		targetText = str
		return !idPattern.MatchString(str)
	})

	var idStartAt int
	for i, rune := range targetText {
		if string(rune) == "0" {
			idStartAt = i
			break
		}
	}

	var (
		idBuilder    strings.Builder
		passedIDlist []string
	)
	position := 0
	for i, rune := range targetText {
		if i >= idStartAt {

			if rune > 'F' /* >'E'> ... >'2'>'1'>'0' */ {
				break
			}

			idBuilder.WriteRune(rune)
			if position == 4 {
				passedIDlist = append(passedIDlist, idBuilder.String())
				idBuilder.Reset()
				position = -1
			}
			position++
		}
	}

	return passedIDlist
}

func convertUTF8toSjis(utf8Str string) string {
	encoder := japanese.ShiftJIS.NewEncoder()
	sjisStr, _, err := transform.String(encoder, utf8Str)
	if err != nil {
		log.Fatal(err)
	}
	return sjisStr
}
