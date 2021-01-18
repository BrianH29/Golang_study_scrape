package scrape

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

//Scrape jobkorea
func Scrape(keyword string) {
	var baseURL string = "http://www.jobkorea.co.kr/Search/?stext=" + keyword

	var jobs []extractedJob
	c := make(chan []extractedJob)

	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
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

func getPage(page int, url string, pageC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)

	pageURL := url + "&tabType=recruit&Page_No=" + strconv.Itoa(page+1)
	fmt.Println("Requesting :", pageURL)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".list-default")

	searchCards.Each(func(i int, cards *goquery.Selection) {
		go extractJob(cards, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	pageC <- jobs
}

func extractJob(cards *goquery.Selection, c chan<- extractedJob) {
	corp := cards.Find(".post-list-corp").Text()
	info := cards.Find(".post-list-info").Text()

	c <- extractedJob{
		corp: CleanString(corp),
		info: CleanString(info)}
}

func getPages(url string) int {
	pages := 0

	res, err := http.Get(url)
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

//CleanString to clear out the sentence
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
