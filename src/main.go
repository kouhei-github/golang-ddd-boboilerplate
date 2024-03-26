package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", helloHandler)
	if err := app.Listen(":8888"); err != nil {
		panic(err)
	}

}

func helloHandler(c *fiber.Ctx) error {
	res := c.Response()
	res.Header.SetStatusCode(fiber.StatusOK)
	return c.SendString("Hello, World 👋!")
}
