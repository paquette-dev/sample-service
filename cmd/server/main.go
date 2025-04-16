package main

import (
	"log"
	"net/http"
	"github.com/labstack/echo/v4"
	"sample-service/internal/database"
)

func main() {
	db, err := database.InitDB("./database.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}