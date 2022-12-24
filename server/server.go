package server

import (
	"log"
	"mindserver/server/mindustryserver"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func Serve() {
	app := fiber.New()
	app.Static("/", "./webpage")

	minServer := mindustryserver.New(app)
	minServer.Route()

	app.Post("/api/post/pull_new_version_restart", func(c *fiber.Ctx) error {
		minServer.Kill()

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
