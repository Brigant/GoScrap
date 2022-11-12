package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Book struct {
	Brand string
	Email string
}

var AllowedDomains = "www.kindundjugend.com"

func main() {
	file, err := os.Create("export.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headers := []string{"Brand", "Email"}
	writer.Write(headers)

	c := colly.NewCollector(
		colly.AllowedDomains(AllowedDomains),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	c.OnHTML("a.slick-next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))

		c.Visit(nextPage)
	})

	book := Book{}

	c.OnHTML(".col.col1ergebnis", func(e *colly.HTMLElement) {

		suffix := e.ChildAttr(".initial_noline", "href")

		book.Brand = e.ChildText(".initial_noline strong")

		suburl := "https://www.kindundjugend.com" + suffix

		c.OnHTML(".sico.ico_email", func(e *colly.HTMLElement) {
			book.Email = e.ChildText(".xsecondarylink span")
		})
		c.Visit(suburl)

		row := []string{book.Brand, book.Email}
		writer.Write(row)
	})

	startUrl := "https://www.kindundjugend.com/kindundjugend-exhibitors/list-of-exhibitors/"
	c.Visit(startUrl)

}