package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"regexp"
	"strings"
)

type Flat struct {
	ID           string
	Price        string
	Sqft_m2      string
	Rooms        string
	Toilets      string
	Area         string
	Elevator     string // No or Yes
	Parking      string // No or Yes
	Heating      string // No or Yes
	CoolAir      string // No or Yes
	RatioEurM    string
	Construction string // new or Not new
	Balcony      string // No or Yes
	Pool         string // No or Yes
	PublicTr     string // No or Yes
	Yard         string // No or Yes
}

var output []string

func main() {

	fName := "pisos_hospitalet.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}

	flat := Flat{}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Price", "Sqft_m2", "RatioEurM", "Rooms", "Toilets", "Area", "Elevator", "Parking", "Heating", "CoolAir", "Construction", "Balcony", "Pool", "PublicTr", "Yard"})

	// Instantiate default collector
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "https://www.habitaclia.com/comprar-") {
			//fmt.Println("the link it has no prefix comprar", link)
			return
		}

		e.Request.Visit(link)
	})

	c.OnHTML(`ul.feature-container`, func(e *colly.HTMLElement) {

		flat.Price = ""
		flat.Sqft_m2 = ""

		parameters := strings.Split(e.ChildText("li.feature"), "\n")

		for _, k := range parameters {
			switch {
			case strings.Contains(k, "€/m2"):
				fmt.Println("\nThis is the RatioEurM: ", k)
				k := strings.TrimSpace(k)
				flat.RatioEurM = k
			case strings.HasSuffix(k, "€"):
				fmt.Println("\nThis is the price:", k)
				flat.Price = k
			case strings.Contains(k, "m2"):
				fmt.Println("\nThis is the Square meters:", k)
				k := strings.TrimSpace(k)
				flat.Sqft_m2 = k
			}
		}

		flat.Balcony = ""
		flat.Parking = ""
		flat.Elevator = ""
		flat.Construction = ""
		flat.Heating = ""
		flat.Toilets = ""
		flat.CoolAir = ""
		flat.Rooms = ""
		flat.Area = ""
		flat.Pool = ""
		flat.PublicTr = ""
		flat.Yard = ""

	})

	c.OnHTML(`article.location`, func(e *colly.HTMLElement) {

		flat.Area = e.ChildText("a.jqVerMapaZonaTooltip")

		fmt.Println("The Area is:", flat.Area)

	})

	c.OnHTML(`section.detail`, func(e *colly.HTMLElement) {

		name := e.Request.URL.Path

		fmt.Println("This is the url:", name)

		re := regexp.MustCompile(`.*-(i.*).htm.*`)

		flat.ID = re.FindStringSubmatch(name)[1]

		// Distribution

		parameters := strings.Split(e.ChildText("article.has-aside"), "\n")

		// fmt.Println(parameters)

		// fmt.Println("this is the ID of the flat\n", flat.ID)

		fmt.Println("the value of balcony is:", flat.Balcony)

		for _, k := range parameters {
			switch {
			case strings.HasSuffix(k, "habitaciones"):
				fmt.Println("\nThis is the number of rooms:", k)
				k := strings.TrimSpace(k)
				flat.Rooms = k
			case strings.HasSuffix(k, "habitación"):
				fmt.Println("\nThis is the number of individual rooms:", k)
				k := strings.TrimSpace(k)
				flat.Rooms = k
			case strings.Contains(k, "Terraza"):
				fmt.Println("The balcony string is:", k)
				fmt.Println("\nBalcony: Yes")
				flat.Balcony = "Yes"
			case strings.HasSuffix(k, "Baños"):
				fmt.Printf("\nThis has %s", k)
				k := strings.TrimSpace(k)
				flat.Toilets = k
			case strings.HasSuffix(k, "Baño"):
				fmt.Printf("\nThis has %s", k)
				k := strings.TrimSpace(k)
				flat.Toilets = k
			case strings.HasSuffix(k, "acondicionado"):
				fmt.Println("\nCoolAir: Yes")
				flat.CoolAir = "Yes"
			case strings.Contains(k, "Sin aire acondicionado"):
				fmt.Println("\nCoolAir: No")
				flat.CoolAir = "No"
			case strings.Contains(k, "Sin plaza parking"):
				fmt.Println("\nParking: No")
				flat.Parking = "No"
			case strings.Contains(k, "Plaza parking"):
				fmt.Println("\nParking: YES")
				flat.Parking = "Yes"
			case strings.Contains(k, "Sin calefacción"):
				fmt.Println("\nHeating: NO")
				flat.Heating = "No"
			case strings.Contains(k, "Calefacción"):
				fmt.Println("\nHeating: Yes")
				flat.Heating = "Yes"
			case strings.Contains(k, "Obra nueva"):
				fmt.Println("\nConstruction: New")
				flat.Construction = "New"
			case strings.Contains(k, "Sin ascensor"):
				fmt.Println("\nAscensor: NO")
				flat.Elevator = "No"
			case strings.Contains(k, "Ascensor"):
				fmt.Println("\nAscensor: YES")
				flat.Elevator = "Yes"
			case strings.HasPrefix(k, "Piscina"):
				fmt.Println("\nPiscina: YES")
				flat.Pool = "Yes"
			case strings.Contains(k, "Cerca de transporte público"):
				fmt.Println("\nTransporte Publico: YES")
				flat.PublicTr = "Yes"
			case strings.Contains(k, "Jardín"):
				fmt.Println("\nJardin: YES")
				flat.Yard = "Yes"
			}
		}

		if flat.CoolAir == "" {
			flat.CoolAir = "No"
		}

		if flat.Elevator == "" {
			fmt.Println("\nAscensor: No")
			flat.Elevator = "No"
		}

		if flat.Heating == "" {
			flat.Heating = "No"
		}

		if flat.Construction == "" {
			flat.Construction = "Not New"
		}

		if flat.Balcony == "" {
			fmt.Println("\nBalcony: NO")
			flat.Balcony = "No"
		}

		if flat.Pool == "" {
			fmt.Println("\nPool: NO")
			flat.Pool = "No"
		}

		if flat.PublicTr == "" {
			fmt.Println("\nPublic Transport: NO")
		    flat.PublicTr = "No"
		}

		if flat.Yard == "" {
			fmt.Println("\nYard: NO")
			flat.Yard = "No"
		}



		output = []string{flat.ID, flat.Price, flat.Sqft_m2, flat.RatioEurM, flat.Rooms, flat.Toilets, flat.Area, flat.Elevator, flat.Parking,
			flat.Heating, flat.CoolAir,  flat.Construction, flat.Balcony, flat.Pool, flat.PublicTr, flat.Yard}

		fmt.Println("escribiendo en output", output)

		writer.Write(output)

	})

	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "\nError:", err)
	})

	c.Visit("https://www.habitaclia.com/viviendas-hospitalet_de_llobregat.htm")
	//c.Visit("https://www.habitaclia.com/viviendas-en-barcelones.htm")

	log.Printf("Scraping finished, check file %q for results\n", fName)

}
