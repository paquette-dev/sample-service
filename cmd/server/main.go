// @title Sample Service API
// @version 1.0
// @description API for managing users
// @host localhost:1323
// @BasePath /
package main

import (
	"log"
	"github.com/labstack/echo/v4"
	"sample-service/internal/database"
	"sample-service/internal/routes"
)

func main() {
	db, err := database.InitDB("./database.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	err = database.SeedDB(db)
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	e := echo.New()
	routes.RegisterUserRoutes(e, db)
	routes.RegisterSwaggerRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}