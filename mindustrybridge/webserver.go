package mindustrybridge

import (
	"mindserver/log"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Server struct {
	mindustry *mindustryBridge
}

type Command struct {
	Cmd string `json:"command" xml:"command" form:"command"`
}

func Route(router *gin.RouterGroup) *Server {
	server := new(Server)
	server.mindustry = newMindustryServer()

	log.Named("Mindustry").Info("Starting server")
	if err := server.mindustry.start(); err != nil {
		log.Panic("Error when starting server: " + err.Error())
	}
	router.POST("/map/:filename", server.uploadNewMap)
	router.GET("/map/:filename", server.downloadMap)
	router.GET("/ws", server.websocketConn)

	return server
}

func (server Server) Kill() {
	server.mindustry.kill()
}

func (server Server) Exit() error {
	return server.mindustry.exit()
}

func (server Server) MustExit() {
	log.Named("Mindustry").Info("Terminating subservice")
	if server.Exit() != nil {
		log.Named("Mindustry").Warn("Cannot terminate subservice, killing")
		server.Kill()
	}
	log.Named("Mindustry").Info("Subservice exited")
}

func (server Server) uploadNewMap(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Error(err.Error())
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer file.Close()
	data := make([]byte, fileHeader.Size)
	file.Read(data)

	saveFile, err := os.Create("mindustrybridge/mindustry-server/config/maps/" + c.Param("filename"))
	if err != nil {
		log.Named("Mindustry").Error("Error when creating file: " + err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	defer saveFile.Close()

	_, err = saveFile.Write(data)
	if err != nil {
		log.Named("Mindustry").Error("Error when writing file: " + err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Named("Mindustry").Info("New map uploaded: " + c.Param("filename"))

	c.Status(http.StatusOK)
}

func (server Server) downloadMap(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func (server Server) websocketConn(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	channel := make(chan []byte)
	server.mindustry.appendOutputChannel(channel)

	// handle websocket closed
	defer func() {
		server.mindustry.RemoveOutputChannel(channel)
		close(channel)
		closeSocketErr := ws.Close()
		if closeSocketErr != nil {
			panic(err)
		}
		log.Named("Mindustry").Info("Websocket connection closed")
	}()

	// send message to client (when message is sent to channel)
	go func() {
		for {
			msg, ok := <-channel
			// close goroutine if channel is being closed
			if !ok {
				runtime.Goexit()
			}
			err := ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Named("Mindustry").Error("Websocket write: " + err.Error())
			}
		}
	}()

	// receive message from client
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Named("Mindustry").Error("Websocket read: " + err.Error())
			break
		}
		log.Named("Mindustry").Info("Websocket recv: " + string(msg))
		server.mindustry.sendCommand(string(msg))
	}
}
