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

// TODO: corregir cuando se pone comillas en el header
// TODO: numeric truncates value
func cleanString(s string) string {
	output := ""
	for _, char := range s {
		if char <= '~' {
			output = fmt.Sprint(output + string(char))
		}
	}
	return output
}

// TODO: make some strategy pattern for these Table, Tables and JSONTable functions
func JSONTable(c echo.Context, db *sql.DB) error {

	tableName := c.Param("tableName")

	tableHeadersQuery := "SELECT column_name FROM information_schema.columns WHERE table_name = '" + tableName + "'"
	tableHeadersRows, err := db.Query(tableHeadersQuery)
	if err != nil {
		panic(err)
	}
	defer tableHeadersRows.Close()

	tableHeaders := make([]string, 0)
	for tableHeadersRows.Next() {
		var header string
		if err := tableHeadersRows.Scan(&header); err != nil {
			log.Fatal(err)
		}
		tableHeaders = append(tableHeaders, header)

	}

	tableValues := make([][]string, 0)
	tableValues = append(tableValues, tableHeaders)

	columns := make([]interface{}, len(tableHeaders))
	columnsPointers := make([]interface{}, len(tableHeaders))
	for idx := 0; idx < len(tableHeaders); idx++ {
		columnsPointers[idx] = &columns[idx]
	}

	tableValuesQuery := "SELECT * FROM " + tableName
	tableValueRow, err := db.Query(tableValuesQuery)
	if err != nil {
		panic(err)
	}
	defer tableValueRow.Close()

	for tableValueRow.Next() {
		if err = tableValueRow.Scan(columnsPointers...); err != nil {
			panic(err)
		}
		values := make([]string, len(columnsPointers))
		for idx, _ := range columnsPointers {
			val := fmt.Sprint(*columnsPointers[idx].(*interface{}))

			values[idx] = val
		}
		tableValues = append(tableValues, values)
	}
	// kdfasjdfksdadjsadkj
	return c.JSON(http.StatusOK, tableValues)
}

func Table(c echo.Context, db *sql.DB) error {
	tableName := c.Param("tableName")

	tableHeadersQuery := "SELECT column_name FROM information_schema.columns WHERE table_name = '" + tableName + "'"
	tableHeadersRows, err := db.Query(tableHeadersQuery)
	if err != nil {
		panic(err)
	}
	defer tableHeadersRows.Close()

	tableHeaders := make([]string, 0)
	for tableHeadersRows.Next() {
		var header string
		if err := tableHeadersRows.Scan(&header); err != nil {
			log.Fatal(err)
		}
		tableHeaders = append(tableHeaders, header)

	}

	tableValues := make([][]string, 0)
	tableValues = append(tableValues, tableHeaders)

	columns := make([]interface{}, len(tableHeaders))
	columnsPointers := make([]interface{}, len(tableHeaders))
	for idx := 0; idx < len(tableHeaders); idx++ {
		columnsPointers[idx] = &columns[idx]
	}

	tableValuesQuery := "SELECT * FROM " + tableName
	tableValueRow, err := db.Query(tableValuesQuery)
	if err != nil {
		panic(err)
	}
	defer tableValueRow.Close()

	for tableValueRow.Next() {
		if err = tableValueRow.Scan(columnsPointers...); err != nil {
			panic(err)
		}
		values := make([]string, len(columnsPointers))
		for idx, _ := range columnsPointers {
			val := fmt.Sprint(*columnsPointers[idx].(*interface{}))

			values[idx] = val
		}
		tableValues = append(tableValues, values)
	}

	return c.Render(http.StatusOK, "table", tableValues)
}

func Tables(c echo.Context, db *sql.DB) error {
	tableNamesQuery := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	tableNamesRows, err := db.Query(tableNamesQuery)
	if err != nil {
		panic(err)
	}
	defer tableNamesRows.Close()

	tableNames := make([]string, 0)

	for tableNamesRows.Next() {
		var table string
		if err := tableNamesRows.Scan(&table); err != nil {
			log.Fatal(err)
		}
		tableNames = append(tableNames, table)
	}
	if err := tableNamesRows.Err(); err != nil {
		panic(err)
	}
	tablesHeaders := make(map[string][]string, 0)
	for _, tableName := range tableNames {
		tableHeadersQuery := "SELECT column_name FROM information_schema.columns WHERE table_name = '" + tableName + "'"
		tableHeadersRows, err := db.Query(tableHeadersQuery)
		if err != nil {
			panic(err)
		}
		defer tableHeadersRows.Close()

		tablesHeaders[tableName] = make([]string, 0)
		for tableHeadersRows.Next() {
			var header string
			if err := tableHeadersRows.Scan(&header); err != nil {
				log.Fatal(err)
			}
			tablesHeaders[tableName] = append(tablesHeaders[tableName], header)
		}
	}
	tableValues := make(map[string][][]string, 0)

	for _, tableName := range tableNames {
		rowsQuery := "SELECT * FROM " + tableName + " LIMIT 10"
		rows, err := db.Query(rowsQuery)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		columns := make([]interface{}, len(tablesHeaders[tableName]))
		columnPointers := make([]interface{}, len(tablesHeaders[tableName]))
		for idx := 0; idx < len(columns); idx++ {
			columnPointers[idx] = &columns[idx]
		}

		tableValues[tableName] = make([][]string, 0)
		tableValues[tableName] = append(tableValues[tableName], tablesHeaders[tableName])
		for rows.Next() {
			rowValues := make([]string, 0)
			err := rows.Scan(columnPointers...)
			if err != nil {
				panic(err)
			}
			for idx := 0; idx < len(columnPointers); idx++ {
				vals := fmt.Sprint(*columnPointers[idx].(*interface{}))
				rowValues = append(rowValues, vals)
			}
			tableValues[tableName] = append(tableValues[tableName], rowValues)

		}

	}

	return c.Render(http.StatusOK, "tables", tableValues)
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
		"decimal": func(value string, amountOfDigits, commaPosition int) bool {
			if commaPosition > 0 {
				if value[len(value)-1-commaPosition] != ',' {
					return false
				}

			} else if commaPosition < 0 {
				if len(value)-commaPosition >= amountOfDigits {
					return false
				}
				for i := len(value) - commaPosition; i < len(value); i++ {
					if value[i] != '0' {
						return false
					}
				}
			} else {
				if len(value) > amountOfDigits {
					return false
				}
			}
			return true
		},
		"integer": func(value string, lowerBound, upperBound int) bool {
			n, err := strconv.Atoi(value)
			if err != nil {
				return false
			}
			return n >= lowerBound && n <= upperBound
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
	for i := 0; i < len(headers); i++ {
		headers[i] = cleanString(headers[i])

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
						stringRow += "'" + columnValue + "'"
					} else {
						rowIsBroken = true
						break
					}

				}
			case "integer":
				{
					lowerBound := c.FormValue("lowerBound_" + header)
					upperBound := c.FormValue("upperBound_" + header)

					lowerBoundNumber, err := strconv.Atoi(lowerBound)
					if err != nil {
						panic(err)
					}

					upperBoundNumber, err := strconv.Atoi(upperBound)
					if err != nil {
						panic(err)
					}
					if typeToCheck[columnType].(func(string, int, int) bool)(columnValue, lowerBoundNumber, upperBoundNumber) {
						stringRow += "'" + columnValue + "'"
					} else {
						rowIsBroken = true
						break
					}

				}
			}
			// TODO handle RowIsBroken
			if i != len(headers)-1 {
				stringRow += ","
			}
		}
		i++
		if !rowIsBroken {
			db.Exec(insertTemplate + stringRow + ")")
		} else {
			fmt.Println("This row is broken: ")
			fmt.Println(row)
		}

	}

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
	e.GET("/tables", func(c echo.Context) error {
		return Tables(c, db)
	})
	e.GET("/table/:tableName", func(c echo.Context) error {
		return Table(c, db)
	})
	e.GET("/table/api/:tableName", func(c echo.Context) error {
		return JSONTable(c, db)
	})

	e.Logger.Fatal(e.Start(":3000"))

}
