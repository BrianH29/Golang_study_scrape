package main

import (
	"os"
	"strings"

	"github.com/BrianH29/Go_scrapper/scrape"
	"github.com/labstack/echo"
)

const fileName string = "jobs.csv"

func home(c echo.Context) error {
	return c.File("home.html")
}

func scrapper(c echo.Context) error {
	defer os.Remove(fileName)
	keyword := strings.ToLower(scrape.CleanString(c.FormValue("keyword")))

	scrape.Scrape(keyword)
	return c.Attachment(fileName, fileName)
}

func main() {
	e := echo.New()
	e.GET("/", home)
	e.POST("/scrape", scrapper)

	e.Logger.Fatal(e.Start(":1323"))
}
