package server

import (
	"fmt"
	"log"
	"mindserver/server/mindustryserver"

	"github.com/gofiber/fiber/v2"
)

func Serve() {
	mindustryServer := mindustryserver.NewMindustryServer()

	app := fiber.New()

	app.Static("/", "./webpage")

	app.Post("/api/post/start_server", func(c *fiber.Ctx) error {
		mindustryServer.Start()
		return nil
	})
	app.Post("/api/post/send_command", func(c *fiber.Ctx) error {
		fmt.Println(c.Body())
		// mindustryServer.SendCommand(string(c.Body()))
		return nil
	})
	app.Get("/api/get/maps_list", func(c *fiber.Ctx) error {
		return nil
	})
	app.Get("/api/get/commandline_output", func(c *fiber.Ctx) error {
		output, err := mindustryServer.GetOutput()
		c.Send(output)
		return err
	})
	app.Post("api/post/runwave", func(c *fiber.Ctx) error {
		return nil
	})

	log.Fatal(app.Listen(":8086"))
}
