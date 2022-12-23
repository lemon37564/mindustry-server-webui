package server

import (
	"log"
	"mindserver/server/mindustryserver"
	"net/http"
	"os"
	"os/exec"

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
	app.Post("/api/post/kill_server", func(c *fiber.Ctx) error {
		if err := mindustryServer.Kill(); err != nil {
			log.Println("Error when killing server:", err)
			return err
		}
		log.Println("[Info] Server killed")
		mindustryServer = mindustryserver.NewMindustryServer()
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
		log.Println("Upload new map:", c.Params("filename"))
		return nil
	})

	app.Post("/api/post/pull_new_version_restart", func(c *fiber.Ctx) error {
		cmd := exec.Command("git", "pull")
		cmd.Run() // wait pull complete
		log.Println("git pull finished, exiting...")
		cmd = exec.Command("go", "run", "mindserver")
		cmd.Start() // create new process and leave
		os.Exit(0)
		panic("unreachable")
	})

	log.Fatal(app.Listen(":8086"))
}
