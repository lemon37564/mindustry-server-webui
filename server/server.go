package server

import (
	"log"
	"mindserver/server/mindustryserver"
	"os"
	"os/exec"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Server struct {
	mindustry *mindustryserver.MindustryServer
	app       *fiber.App
}

type Command struct {
	Cmd string `json:"command" xml:"command" form:"command"`
}

func New() Server {
	server := *new(Server)
	server.mindustry = mindustryserver.NewMindustryServer()
	server.app = fiber.New()
	return server
}

func (server Server) Serve() {
	server.hanleSigInt()
	server.app.Static("/", "./webpage")

	server.app.Post("/api/post/start_server", func(c *fiber.Ctx) error {
		if err := server.mindustry.Start(); err != nil {
			log.Println("Error when starting server:", err)
			return err
		}
		return nil
	})
	server.app.Post("/api/post/kill_server", func(c *fiber.Ctx) error {
		if err := server.mindustry.Kill(); err != nil {
			log.Println("Error when killing server:", err)
			return err
		}
		log.Println("[Info] Server killed")
		server.mindustry = mindustryserver.NewMindustryServer()
		return nil
	})
	server.app.Post("/api/post/upload_new_map/:filename", func(c *fiber.Ctx) error {
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

	server.app.Post("/api/post/pull_new_version_restart", func(c *fiber.Ctx) error {
		server.mindustry.Kill()

		cmd := exec.Command("git", "pull")
		cmd.Run() // wait pull complete
		log.Println("git pull finished, exiting...")
		cmd = exec.Command("go", "run", "mindserver")
		cmd.Start() // create new process and leave

		os.Exit(0)
		panic("unreachable")
	})

	server.app.Get("/ws/mindustry_server", websocket.New(func(c *websocket.Conn) {
		closed := false
		// handle websocket closed
		// should have better way to handle this
		c.SetCloseHandler(func(code int, text string) error {
			closed = true
			if code == websocket.CloseNormalClosure {
				log.Println("Websocket connection closed")
			}
			return nil
		})
		go func() {
			for {
				msg := <-server.mindustry.GetOutputChannel()
				// close goroutine if connection closed
				if closed {
					return
				}
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Websocket write:", err)
				}
			}
		}()
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Websocket read:", err)
				break
			}
			log.Printf("Websocket recv: %s", msg)
			server.mindustry.SendCommand(string(msg))
		}
	}))

	log.Fatal(server.app.Listen(":8086"))
}

func (server Server) hanleSigInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		sig := <-c
		log.Println("Recived signal", sig)
		log.Println("Killing mindustry server")
		server.mindustry.Kill()
		log.Println("Exiting")
		os.Exit(0)
	}()
}
