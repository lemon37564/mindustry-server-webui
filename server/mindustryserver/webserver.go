package mindustryserver

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Server struct {
	mindustry *mindustryBridge
	app       *fiber.App
}

type Command struct {
	Cmd string `json:"command" xml:"command" form:"command"`
}

func New(app *fiber.App) *Server {
	server := new(Server)
	server.mindustry = newMindustryServer()
	if err := server.mindustry.start(); err != nil {
		log.Fatal("Error when starting server:", err)
	}
	server.app = app
	return server
}

func (server Server) Route() {
	server.handleSigInt()

	server.app.Post("/api/post/force_restart_server", server.forceRestartServer)
	server.app.Post("/api/post/upload_new_map/:filename", server.uploadNewMap)

	server.app.Get("/ws/mindustry_server", websocket.New(server.websocketConn))
}

func (server Server) Kill() {
	server.mindustry.kill()
}

func (server Server) forceRestartServer(c *fiber.Ctx) error {
	if err := server.mindustry.kill(); err != nil {
		log.Println("Error when killing server:", err)
		return err
	}
	log.Println("[Info] Server killed, restarting")
	server.mindustry = newMindustryServer()
	if err := server.mindustry.start(); err != nil {
		log.Fatal("Error when starting server:", err)
	}
	return nil
}

func (server Server) uploadNewMap(c *fiber.Ctx) error {
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
}

func (server Server) websocketConn(c *websocket.Conn) {
	channel := make(chan []byte)
	server.mindustry.appendOutputChannel(channel)
	// handle websocket closed, close the channel
	c.SetCloseHandler(func(code int, text string) error {
		server.mindustry.RemoveOutputChannel(channel)
		close(channel)
		if code == websocket.CloseNormalClosure {
			log.Println("Websocket connection closed")
		}
		return nil
	})
	// send message to client (when message is sent to channel)
	go func() {
		for {
			msg, ok := <-channel
			// close goroutine if channel is being closed
			if !ok {
				return
			}
			err := c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Websocket write:", err)
			}
		}
	}()
	// receive message from client
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Websocket read:", err)
			break
		}
		log.Printf("Websocket recv: %s", msg)
		server.mindustry.sendCommand(string(msg))
	}
}

func (server Server) handleSigInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		sig := <-c
		log.Println("Recived signal", sig)
		log.Println("Killing mindustry server")
		server.mindustry.kill()
		log.Println("Exiting")
		os.Exit(0)
	}()
}
