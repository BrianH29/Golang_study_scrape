package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	corp string
	info string
}

var baseURL string = "http://www.jobkorea.co.kr/Search/?stext=golang"

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	createFile(jobs)
	fmt.Println("Finish, extraction", len(jobs))
}

func createFile(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	header := []string{"Corporation", "job info"}
	hErr := w.Write(header)
	checkErr(hErr)

	for _, job := range jobs {
		jobSlice := []string{job.corp, job.info}
		jErr := w.Write(jobSlice)
		checkErr(jErr)
	}

}

func getPage(page int) []extractedJob {
	var jobs []extractedJob

	pageURL := baseURL + "&tabType=recruit&Page_No=" + strconv.Itoa(page+1)
	fmt.Println("Requesting :", pageURL)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".list-default")

	searchCards.Each(func(i int, cards *goquery.Selection) {
		job := extractJob(cards)
		jobs = append(jobs, job)
	})

	return jobs
}

func extractJob(cards *goquery.Selection) extractedJob {
	corp := cards.Find(".post-list-corp").Text()
	info := cards.Find(".post-list-info").Text()

	return extractedJob{
		corp: cleanString(corp),
		info: cleanString(info)}
}

func getPages() int {
	pages := 0

	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".tplPagination.wide").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("li").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
