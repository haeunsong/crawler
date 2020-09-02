package contacts

import (
	"crawler/tools"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	_ "github.com/mattn/go-sqlite3"

	"github.com/PuerkitoBio/goquery"
)

var id2 = 0

type DepartmentOption struct {
	id       int
	name     string
	major    string
	location string // 연구실 위치
	number   string
	email    string
}

func addData(extractedNumbers []DepartmentOption, targetURL2 string, numbers []DepartmentOption, index int, index2 int) []DepartmentOption {
	extractedNumbers = getDepartmentPage(targetURL2, index, index2)
	numbers = append(numbers, extractedNumbers...)

	return numbers
}
func GetDepartmentData() {
	var numbers []DepartmentOption
	var extractedNumbers []DepartmentOption
	// 신학과
	// 일어일본학과
	// 중어중국학과
	var targetURL = [20]string{"http://sinhak.skhu.ac.kr/icons/app/cms/?html=/home/prof.html&shell=/index.shell:15",
		"http://depjs.skhu.ac.kr/icons/app/cms/?html=/home/prof.html&shell=/index.shell:15",
		"http://depcs.skhu.ac.kr/icons/app/cms/?html=/home/edu_course_ilban.html&shell=/index.shell:56",
	}
	for i := 0; i < 3; i++ {
		extractedNumbers = getDepartmentPage(targetURL[i], 0, 0)
		numbers = append(numbers, extractedNumbers...)
	}
	// 사회복지학과(ok) 수정 필요
	targetURL2 := "http://welfare.skhu.ac.kr/icons/app/cms/?html=/home/future.html&shell=/index.shell:42"
	numbers = addData(extractedNumbers, targetURL2, numbers, 1, 0)
	// 사회과학부(ok) 수정 필요
	targetURL2 = "http://sscience.skhu.ac.kr/icons/app/cms/?html=/home/int_manag1.html&shell=/index.shell:24"
	numbers = addData(extractedNumbers, targetURL2, numbers, 2, 1)
	// 신문방송학과(ok) 수정 필요
	targetURL2 = "http://media.skhu.ac.kr/icons/app/cms/?html=/home/int1_2_1.html&shell=/index.shell:195"
	numbers = addData(extractedNumbers, targetURL2, numbers, 3, 2)
	// 경영학부
	targetURL2 = "http://biz.skhu.ac.kr/"
	numbers = addData(extractedNumbers, targetURL2, numbers, 0, 0)
	// 컴공(ok)
	targetURL2 = "http://cse.skhu.ac.kr/icons/app/cms/?html=/home/int1_3.html&shell=/index.shell:26"
	numbers = addData(extractedNumbers, targetURL2, numbers, 4, 3)
	// 소프(ok)
	targetURL2 = "http://sw.skhu.ac.kr/icons/app/cms/?html=/home/int1_2.html&shell=/index.shell:132"
	numbers = addData(extractedNumbers, targetURL2, numbers, 5, 4)
	// 정통
	targetURL2 = "http://cc.skhu.ac.kr/icons/app/cms/?html=/home/int3_1.html&shell=/index.shell:229"
	numbers = addData(extractedNumbers, targetURL2, numbers, 6, 5)
	// 디컨
	targetURL2 = "http://dicon.skhu.ac.kr/sub/sub0201.php"
	numbers = addData(extractedNumbers, targetURL2, numbers, 7, 6)

	// WriteNumbers2(numbers)
	fmt.Println("Done, extracted", len(numbers))

}

func getDepartmentPage(targetURL string, index int, index2 int) []DepartmentOption {
	check := []string{"table.table > tbody ", "div.prof > table > tbody", "div.all > .contents", "div.proboxB > ul",
		"div.box > ul", "table.tbTypeA > tbody", "table.B_type > tbody",
		"div.professor_desc.fr"}

	var numbers []DepartmentOption
	fmt.Println("Requesting ", targetURL)
	res, err := http.Get(targetURL)

	fmt.Println(res)

	tools.CheckErr(err)
	tools.CheckCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	//doc, err := goquery.NewDocumentFromReader(tools.EucKrReaderToUtf8Reader(res.Body))
	tools.CheckErr(err)

	doc.Find(check[index]).Each(func(i int, s *goquery.Selection) {
		number := extractDepartmentNumber(s, index2)
		WriteNumbers2(number.id, number.name, number.major, number.location, number.number, number.email)
		numbers = append(numbers, number)
	})

	return numbers
}

func extractDepartmentNumber(s *goquery.Selection, index2 int) DepartmentOption {
	// var id int
	id2++
	var name, major, location, number, email string
	switch index2 {
	case 0:
		id = id2
		name = strings.Trim(strings.Trim(tools.CleanString(s.Children().Eq(0).Find("td").Text()), "· "), "교수")
		major = strings.Trim(tools.CleanString(s.Children().Eq(1).Find("td").Text()), "· ")
		location = strings.Trim(tools.CleanString(s.Children().Eq(2).Find("td").Text()), "· ")
		number = strings.Trim(tools.CleanString(s.Children().Eq(3).Find("td").Text()), "· ")
		email = strings.Trim(tools.CleanString(s.Children().Eq(4).Find("td").Text()), "· ")
		break
	case 1: // 사회과학부
		id = id2
		name = strings.Trim(tools.CleanString(s.Children().Eq(0).Text()), "ㆍ이 름 :")
		//major = "사회과학부"
		major = strings.Trim(tools.CleanString(s.Children().Eq(1).Text()), "ㆍ전 공: ")
		location = strings.Trim(tools.CleanString(s.Children().Eq(2).Text()), "ㆍ연구실 : ")
		number = strings.TrimFunc(tools.CleanString(s.Children().Eq(3).Text()), func(r rune) bool {
			return !unicode.IsNumber(r)
		})
		email = strings.Trim(tools.CleanString(s.Children().Eq(4).Text()), "ㆍE-mail :")
		break
	case 2: // 신문방송학과
		id = id2
		name = strings.Trim(tools.CleanString(s.Find(".pro_name").Not("span").Text()), "교수")
		major = "신문방송학과"
		// major = tools.CleanString(s.Children().Eq(1).Not("b").Text())
		location = strings.Trim(tools.CleanString(s.Children().Eq(2).Text()), "연구실 :")
		number = strings.TrimFunc(tools.CleanString(s.Children().Eq(3).Text()), func(r rune) bool {
			return !unicode.IsNumber(r)
		})
		email = strings.Trim(tools.CleanString(s.Children().Eq(4).Text()), "E-mail :")
		break
	case 3: // 컴공 (정연식 교수 수정 필요)
		id = id2
		name = tools.CleanString(s.Children().Eq(0).Text())
		major = "컴퓨터공학과"
		// major = tools.CleanString(s.Children().Eq(7).Text())
		location = strings.Trim(tools.CleanString(s.Children().Eq(2).Text()), "연구실 :")
		number = strings.TrimFunc(tools.CleanString(s.Children().Eq(3).Text()), func(r rune) bool {
			return !unicode.IsNumber(r)
		})
		email = tools.CleanString(s.Children().Eq(4).Find("a").Text())
		break
	case 4: // 소프
		id = id2
		name = tools.CleanString(s.Children().Eq(0).Find("td").Text())
		major = "소프트웨어공학과"
		location = tools.CleanString(s.Children().Eq(1).Find("td").Eq(0).Text())
		number = tools.CleanString(s.Children().Eq(1).Find("td").Eq(1).Text())
		email = tools.CleanString(s.Children().Eq(2).Find("td > a").Text())
		break
	case 5: // 정통
		id = id2
		name = tools.CleanString(s.Children().Eq(0).Find("td > span").Text())
		major = "정보통신공학과"
		//major = tools.CleanString(s.Children().Eq(3).Find("td").Eq(2).Text())
		location = " "
		number = " "
		email = tools.CleanString(s.Children().Eq(1).Find("td > a").Text())
		break
	case 6: // 디컨
		id = id2
		name = tools.CleanString(s.Children().Find("h3").Text()) // 이름 안나옴..
		major = "디지털컨텐츠학과"
		//major = tools.CleanString(s.Children().Find("dl > dd.vertical").Eq(2).Text())
		location = " "
		number = " "
		email = tools.CleanString(s.Children().Find("span.h3_email").Text())
		break

	}

	return DepartmentOption{
		id:       id2,
		name:     name,
		major:    major,
		location: location,
		number:   number,
		email:    email,
	}
}
func WriteNumbers2(id int, name, major, location, number, email string) {
	// sql.DB 객체 생성
	db, err := sql.Open("sqlite3", "./data2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS pages (
			id integer PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(100),
			major VARCHAR(100),
			location VARCHAR(100),
			number VARCHAR(100),
			email VARCHAR(100)
		)`)
	statement.Exec()
	rows, err := db.Query("INSERT INTO pages VALUES (?,?,?,?,?,?)", id, name, major, location, number, email)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &major, &location, &number, &email)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("finish!")

}
