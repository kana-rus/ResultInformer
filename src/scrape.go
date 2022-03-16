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
			isUTF8site = (charset == "utf-8")
		}
		return !exists
		// break if !exists is false, in other words, exists is true.
	})

	if isUTF8site {
		targetCode = targetWord
	} else {
		targetCode = convertUTF8toSjis(targetWord)
	}

	doc.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		href, exists := s.Attr("href")
		isTargetTag := (s.Text() == targetCode)

		if exists && isTargetTag {
			targetHref = href
		}
		return !isTargetTag
		// break if !isTagetTag is false, in other words, isTargetTag is true.
	})

	return targetHref
}


// old version
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


/*
// new version
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
		// 返り値が false になったとき終了
		// つまり str が idPattern にマッチしたら終了
		// どうせ１つ (巨大な塊) しかないのでこれでよし
	})

	var (
		idBuilder    strings.Builder
		passedIDlist []string
		idStart      int
	)
	// const idLen = 5

	// ID がどこからはじまるか調べる
	// ごく最初の方から始まるはずなので、単純に前から見ていけばOK
	for i, rune := range targetText {
		if string(rune) == "0" {
			idStart = i
			break
		}
	}

	// ID の何文字目 (0,1,2,3,4) か調べて繋げる
	// 5 文字繋がったら passedIDlist に入れる
	position := 0
	for i, rune := range targetText {
		str := string(rune)

		if i >= idStart {

			// 終了条件 (ここからは ID じゃない)
			// (「以上〜名」という漢字で引っかかる)
			// サイトの文字コードが ShiftJIS のため面倒
			if rune > 'F' /* >'E'> ... >'2'>'1'>'0' // {
				break
			}

			idBuilder.WriteString(str)
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
*/

func convertUTF8toSjis(utf8Str string) string {
	encoder := japanese.ShiftJIS.NewEncoder()
	sjisStr, _, err := transform.String(encoder, utf8Str)
	if err != nil {
		log.Fatal(err)
	}
	return sjisStr
}
