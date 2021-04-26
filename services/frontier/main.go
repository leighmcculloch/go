package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const archiveURL = "https://history.stellar.org/prd/core-live/core_live_001/"

func main() {
	data, err := NewData(archiveURL)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			data.Update()
			time.Sleep(5 * time.Second)
		}
	}()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := http.StatusText(http.StatusInternalServerError)
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}
			return ctx.Status(code).JSON(struct {
				StatusCode int    `json:"status_code"`
				Error      string `json:"error"`
			}{
				StatusCode: code,
				Error:      message,
			})
		},
	})
	app.Use(recover.New())
	app.Get("/", (&RootHandler{
		ArchiveURL: archiveURL,
		Data:       data,
	}).Handler)
	app.Get("/accounts/:id", (&AccountHandler{
		Data: data,
	}).Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen(":" + port))
}
