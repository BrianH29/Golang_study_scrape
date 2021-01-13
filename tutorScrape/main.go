package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/gommon/log"
)

var baseURL string = "https://www.iei.or.kr/intro/teacher.kh"

type tutor struct {
	name     string
	info     string
	subInfo  string
	intro    string
	subIntro string
}

//Scrape kh tutuor list
func main() {
	var teachers []tutor
	totalList := getPage()

	for _, list := range totalList {
		teachers = append(teachers, list)
	}

	writeTutor(teachers)
	fmt.Println("Done, extraction", len(teachers))
}

func writeTutor(teachers []tutor) {
	file, err := os.Create("tutorList.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	header := []string{"Tutor", "info", "subInfo", "intro", "subIntro"}

	wErr := w.Write(header)
	checkErr(wErr)

	for _, tutor := range teachers {
		tutorSlice := []string{tutor.name, tutor.info, tutor.subInfo, tutor.intro, tutor.subIntro}
		twErr := w.Write(tutorSlice)
		checkErr(twErr)
	}
}

func getPage() []tutor {
	var teachers []tutor
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	tutorList := doc.Find(".intro_list li")

	tutorList.Each(func(i int, section *goquery.Selection) {
		teacher := extractTutor(section)
		teachers = append(teachers, teacher)
	})

	return teachers
}

func extractTutor(section *goquery.Selection) tutor {
	name := section.Find(".intro_name").Text()
	info := section.Find(".intro_info").Text()
	subInfo := section.Find(".intro_answer").Text()
	intro := section.Find(".intro_content1").Text()
	subIntro := section.Find(".intro_content2").Text()

	return tutor{
		name:     cleanString(name),
		info:     cleanString(info),
		subInfo:  cleanString(subInfo),
		intro:    cleanString(intro),
		subIntro: cleanString(subIntro)}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("status code error : %d %s", res.StatusCode, res.Status)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
