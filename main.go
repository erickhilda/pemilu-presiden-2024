package main

import (
	"encoding/csv"
	"fmt"
	"html"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

const url = "https://en.wikipedia.org/wiki/2024_Indonesian_presidential_election"

func main() {
	client := &http.Client{}

	resp, err := client.Get(url)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	file, err := os.Create("data/pemilu-2024.csv")

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	data := make([][]string, 0)

	tableDocs := doc.Find("table.wikitable.sortable")
	tableDocs.Each(func(i int, table *goquery.Selection) {
		if i == 1 {
			table.Find("tr").Each(func(i int, tr *goquery.Selection) {
				row := make([]string, 0)
				tr.Find("th").Each(func(j int, th *goquery.Selection) {
					row = append(row, html.UnescapeString(th.Text()))
				})
				tr.Find("td").Each(func(k int, td *goquery.Selection) {
					row = append(row, html.EscapeString(td.Text()))
				})

				data = append(data, row)
			})
		}
	})

	data = data[4:]
	for i := 0; i < len(data); i++ {
		if len(data[i]) > 8 {
			data[i] = data[i][1:]
		}
	}
	headers := []string{"province", "total", "perc", "votes", "perc", "votes", "perc", "total"}

	data = append([][]string{headers}, data...)

	writer.WriteAll(data)
}
