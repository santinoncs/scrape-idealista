package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"strings"
)

type Flat struct {
	Name        string
	Price       int
	Sqft_m2     int
	Rooms       int
	Toilets		int
	Area        string
	Elevator    string
}

func main() {

	fName := "pisos.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Name", "Price (EUR)", "Sqft_m2", "Rooms", "Area", "Elevator"})

	// Instantiate default collector
	c := colly.NewCollector(
		//colly.AllowedDomains("idealista.com", "www.idealista.com"),
	)



	// details-property_features class

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "https://www.habitaclia.com/comprar") {
			//fmt.Println("the link it has no prefix comprar", link)
			return
		}
		// Print link
		// fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		//e.Request.Visit(e.Request.AbsoluteURL(link))
		e.Request.Visit(link)
	})

	c.OnHTML(".feature-container", func(e *colly.HTMLElement) {
/*		if e.DOM.Find("section.course-info").Length() == 0 {
			return
		}*/

		fmt.Println("class feature found", e.Request.URL)

		price := strings.Split(e.ChildText(".feature"), "\n")[0]
		sqft_m2 := strings.Split(e.ChildText(".feature"), "\n")[1]
		rooms := strings.Split(e.ChildText(".feature"), "\n")[2]
		toilets := strings.Split(e.ChildText(".feature"), "\n")[3]

		fmt.Println("this is the price", price)
		fmt.Println("this is the sqft m2", sqft_m2)
		fmt.Println("this is the sqft m2", rooms)
		fmt.Println("this is the sqft m2", toilets)


	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "\nError:", err)
	})

	//c.Visit("https://www.idealista.com/venta-viviendas/barcelona/eixample/l-antiga-esquerra-de-l-eixample/")
	c.Visit("https://www.habitaclia.com/viviendas-en-barcelones.htm")

	log.Printf("Scraping finished, check file %q for results\n", fName)



}
