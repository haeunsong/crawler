package contacts

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"crawler/tools"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

type InstutitionOption struct {
	id           int
	belong       string
	name         string
	post         string
	task         string
	schoolNumber string
	faxNumber    string
}

const SkhuURL string = "http://skhu.ac.kr"

var targetURL string
var nums = [5]string{"2", "9", "4", "5", "7"}
var id = 0

func GetInstutionData() {
	var numbers []InstutitionOption

	for i := 0; i < 5; i++ {
		targetURL = fmt.Sprintf("%s/uni_int/uni_int_5_%s.aspx", SkhuURL, nums[i])
		extractedNumbers := getInstutionPage(targetURL)
		numbers = append(numbers, extractedNumbers...)
	}

}

func getInstutionPage(targetURL string) []InstutitionOption {

	var numbers []InstutitionOption
	fmt.Println("Requesting ", targetURL)
	res, err := http.Get(targetURL)

	fmt.Println(res)

	tools.CheckErr(err)
	tools.CheckCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(tools.EucKrReaderToUtf8Reader(res.Body))
	tools.CheckErr(err)

	searchNums := doc.Find("table.cont_a.mt20.ml20.w690 > tbody > tr")

	searchNums.Each(func(i int, s *goquery.Selection) {
		number := extractInstuitionNumber(s)
		WriteNumbers(number.id, number.belong, number.name, number.post, number.task, number.schoolNumber, number.faxNumber)
		numbers = append(numbers, number)
	})

	return numbers
}

func extractInstuitionNumber(s *goquery.Selection) InstutitionOption {
	id++
	belong := tools.CleanString(s.Children().Eq(0).Text())
	name := tools.CleanString(s.Children().Eq(1).Text())
	post := tools.CleanString(s.Children().Eq(2).Text())
	task := tools.CleanString(s.Children().Eq(3).Text())
	schoolNumber := tools.CleanString(s.Children().Eq(4).Text())
	faxNumber := tools.CleanString(s.Children().Eq(5).Text())

	return InstutitionOption{
		id:           id,
		belong:       belong,
		name:         name,
		post:         post,
		task:         task,
		schoolNumber: schoolNumber,
		faxNumber:    faxNumber,
	}
}

// id, belong 등등 사용..
func WriteNumbers(id int, belong, pname, post, task, schoolNum, faxNum string) {

	// sql.DB 객체 생성
	db, err := sql.Open("sqlite3", "./data1.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS pages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			belong VARCHAR(100),
			pname VARCHAR(100),
			post VARCHAR(100),
			task VARCHAR(100),
			schoolNum VARCHAR(100),
			faxNum VARCHAR(100)
		  )`)
	statement.Exec()
	rows, err := db.Query("INSERT INTO pages VALUES (?,?,?,?,?,?,?)", id, belong, pname, post, task, schoolNum, faxNum)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &belong, &pname, &post, &task, &schoolNum, &faxNum)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("finish!")

}
