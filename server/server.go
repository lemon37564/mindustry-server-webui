package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func startMindustryServer() (cmd *exec.Cmd, inPipe io.WriteCloser, outPipe *bufio.Reader) {
	serverCmd := exec.Command("java", "-jar", "./mindustry-server/server.jar")

	procIn, _ := serverCmd.StdinPipe()
	procOut, _ := serverCmd.StdoutPipe()
	err := serverCmd.Start()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(procOut)

	procIn.Write([]byte("host\n"))
	output, _, err := reader.ReadLine()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}

func Serve() {
	startMindustryServer()

	app := fiber.New()

	app.Static("/", "./webpage")

	app.Get("/api/get/maps_list", func(c *fiber.Ctx) error {
		return nil
	})
	app.Post("api/post/runwave", func(c *fiber.Ctx) error {
		return nil
	})

	log.Fatal(app.Listen(":8086"))
}
