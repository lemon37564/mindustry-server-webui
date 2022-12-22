package server

import (
	"log"
	"mindserver/server/mindustryserver"

	"github.com/gofiber/fiber/v2"
)

type Command struct {
	Cmd string `json:"command" xml:"command" form:"command"`
}

func Serve() {
	mindustryServer := mindustryserver.NewMindustryServer()

	app := fiber.New()

	app.Static("/", "./webpage")

	app.Post("/api/post/start_server", func(c *fiber.Ctx) error {
		return mindustryServer.Start()
	})
	app.Post("/api/post/send_command", func(c *fiber.Ctx) error {
		cmd := new(Command)
		if err := c.BodyParser(cmd); err != nil {
			log.Println("In parsing body:", err)
			return err
		}
		mindustryServer.SendCommand(cmd.Cmd)
		return nil
	})
	app.Get("/api/get/commandline_output", func(c *fiber.Ctx) error {
		output := mindustryServer.GetOutput()
		return c.Send(output)
	})

	log.Fatal(app.Listen(":8086"))
}
