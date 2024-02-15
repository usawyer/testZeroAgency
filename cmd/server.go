package main

import (
	"github.com/gofiber/fiber/v2"
	database "github.com/usawyer/testZeroAgency/internal/db"
	"github.com/usawyer/testZeroAgency/internal/handlers"
	"github.com/usawyer/testZeroAgency/internal/service"
	"log"
)

func main() {
	db := database.New()
	srvc := service.New(db)
	handler := handlers.New(srvc)
	app := fiber.New()

	handler.InitRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
