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
	query := "CREATE TABLE " + tableName + "("
	for i, header := range headers {
		header = strings.TrimSpace(header)
		fmt.Println(header)
		columnType := c.FormValue("select" + header)
		columnName := c.FormValue("input" + header)
		fmt.Println(columnName + ": " + columnType)
		query += columnName + " " + columnType
		if i != len(headers)-1 {
			query += ","
		}
	}
	query += ")"

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
