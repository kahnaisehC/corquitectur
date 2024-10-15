package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

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

func upload(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	values := make(map[string]string)
	values["Name"] = name
	values["Email"] = email
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

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
	e.POST("/upload", upload)
	e.Logger.Fatal(e.Start(":3000"))

}
