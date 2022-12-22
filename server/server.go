package server

import (
	"log"
	"mindserver/server/mindustryserver"
	"net/http"
	"os"

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
		if err := mindustryServer.Start(); err != nil {
			log.Println("Error when starting server:", err)
			return err
		}
		return nil
	})
	app.Post("/api/post/send_command", func(c *fiber.Ctx) error {
		cmd := new(Command)
		if err := c.BodyParser(cmd); err != nil {
			log.Println("Error in parsing body:", err)
			return err
		}
		if err := mindustryServer.SendCommand(cmd.Cmd); err != nil {
			log.Println("Error in sending command:", err)
			return err
		}
		return nil
	})
	app.Get("/api/get/commandline_output", func(c *fiber.Ctx) error {
		// update when output updated or the force_update is set to true
		if mindustryServer.IsOutputUpdated() || c.Query("force_update") == "true" {
			output := mindustryServer.GetOutput()
			return c.SendString(string(output))
		}
		// the result is not changed, use the cached result
		return c.SendStatus(http.StatusNotModified)
	})
	app.Post("/api/post/upload_new_map/:filename", func(c *fiber.Ctx) error {
		c.Accepts("application/octet-stream")
		f, err := os.Create("config/maps/" + c.Params("filename"))
		if err != nil {
			log.Println("Error when creating file:", err)
			return err
		}
		_, err = f.Write(c.Body())
		if err != nil {
			log.Println("Error when writing file:", err)
			return err
		}
		return nil
	})

	log.Fatal(app.Listen(":8086"))
}
