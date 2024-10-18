package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

func SuccessfullUpload(c echo.Context, vals *map[string]string) error {
	return c.Render(http.StatusOK, "successfullUpload", vals)
}

func upload(c echo.Context, db *sql.DB) error {
	var typeToCheck = map[string]any{
		"varchar": func(varchar string, length int) bool {
			return len(varchar) <= length
		},
		"decimal": func(amountOfDigits, commaPosition int) bool {
			return true
		},
		"integer": func(lowerBound, upperBound int) bool {
			return true
		},
	}

	values := make(map[string]string)

	file, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}
	src, err := file.Open()

	if err != nil {
		panic(err)
	}

	defer src.Close()

	csvFile := csv.NewReader(src)
	headers, err := csvFile.Read()
	if err != nil {
		panic(err)
	}
	tableName := strings.TrimSuffix(file.Filename, ".csv")
	query := "CREATE TABLE IF NOT EXISTS " + tableName + "("
	insertTemplate := "INSERT INTO " + tableName + "("
	for i, header := range headers {
		header = strings.TrimSpace(header)
		columnName := c.FormValue("columnName_" + header)
		columnType := c.FormValue("columnType_" + header)

		query += columnName + " " + columnType
		insertTemplate += columnName
		if i != len(headers)-1 {
			insertTemplate += ", "
			query += ","
		}
	}
	insertTemplate += ") VALUES("
	query += ")"
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	var i int = 0
	for row, err := csvFile.Read(); err != io.EOF; row, err = csvFile.Read() {
		if err != nil {
			panic(err)
		}

		if len(row) != len(headers) {
			continue
		}
		stringRow := ""
		rowIsBroken := false
		for i := 0; i < len(headers); i++ {
			header := strings.TrimSpace(headers[i])
			columnValue := strings.TrimSpace(row[i])
			columnType := c.FormValue("columnType_" + header)
			switch columnType {
			case "varchar":
				{
					varcharLength := c.FormValue("varcharLength_" + header)
					length, err := strconv.Atoi(varcharLength)
					if err != nil {
						panic(err)
					}
					if typeToCheck[columnType].(func(string, int) bool)(columnValue, length) {
						stringRow += "'" + columnValue + "'"
					} else {
						rowIsBroken = true
						break

					}
				}
			case "decimal":
				{
					amountOfDigits := c.FormValue("amountOfDigits_" + header)
					commaPosition := c.FormValue("commaPosition_" + header)

					digits, err := strconv.Atoi(amountOfDigits)
					if err != nil {
						panic(err)
					}
					comma, err := strconv.Atoi(commaPosition)
					if err != nil {
						panic(err)
					}
					if typeToCheck[columnType].(func(string, int, int) bool)(columnValue, digits, comma) {

					}

				}
			case "integer":
				{

				}
			}
			if i != len(headers)-1 {
				stringRow += ","
			}
		}
		i++
		if !rowIsBroken {
			db.Exec(insertTemplate + stringRow + ")")
		}

	}
	fmt.Println(query)

	// create table with headers

	// _, err = db.Exec(query)
	// if err != nil {
	// 	panic(err)
	// }
	return SuccessfullUpload(c, &values)

}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	USER := os.Getenv("USERNAME")
	PASSWORD := os.Getenv("SECRET")
	DB := os.Getenv("DB")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", USER, PASSWORD, DB, HOST, PORT)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection successfull!1!")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "public")
	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}

	e.Renderer = t

	e.GET("/", Index)
	e.POST("/upload", func(c echo.Context) error {
		return upload(c, db)
	})
	e.Logger.Fatal(e.Start(":3000"))

}
