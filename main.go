package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	// "io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	// "os"
	"database/sql"
	"strconv"
	"strings"
)

func main() {
	db, err := sqlx.Connect("mysql", "root:passw0rd@(localhost:3306)/bookszilla?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	db.Ping()

	for i := 1; i <= 1000; i++ {
		url := "http://wmate.ru/ebooks/book" + strconv.Itoa(i) + ".html"
		response, err := http.Get(url)
		println(url + " " + response.Status)

		if err != nil {
			fmt.Printf("%s", err)
			// os.Exit(1)
		} else if response.StatusCode == 200 {
			defer response.Body.Close()

			doc, err := goquery.NewDocumentFromResponse(response)
			if err != nil {
				log.Fatal(err)
			}

			doc.Find("article.post").Each(func(i int, s *goquery.Selection) {
				title := s.Find("h1").Text()
				println(title)

				// author := s.Find("ul.info li").Eq(2).Find("em").Text()
				// println(author)

				var desctiption sql.NullString
				desc := ""
				s.Find("p").Each(func(j int, sp *goquery.Selection) {
					desc += strings.TrimSpace(sp.Text()) + "\n"
				})

				if len(desc) > 0 {
					desctiption.Scan(desc)
				}

				// println(desctiption)

				sql := "INSERT INTO `books` (`name`, `description`) VALUES (?, ?);"
				db.Exec(sql, title, desctiption)
			})
		}
	}
}
