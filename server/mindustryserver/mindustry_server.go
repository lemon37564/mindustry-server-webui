package mindustryserver

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"regexp"
)

type MindustryServer struct {
	cmd           *exec.Cmd
	started       bool
	inPipe        io.WriteCloser
	outPipe       io.ReadCloser
	reader        *bufio.Reader
	outputBuffer  []byte
	outputChanged bool
}

func NewMindustryServer() MindustryServer {
	server := *new(MindustryServer)

	server.cmd = exec.Command("java", "-jar", "./mindustry-server/server.jar")
	server.started = false
	server.inPipe, _ = server.cmd.StdinPipe()
	server.outPipe, _ = server.cmd.StdoutPipe()
	server.outputBuffer = make([]byte, 0)
	server.outputChanged = false
	return server
}

func (server *MindustryServer) Start() (err error) {
	if server.started {
		return
	}
	err = server.cmd.Start()
	if err != nil {
		return err
	}

	server.reader = bufio.NewReader(server.outPipe)
	server.started = true

	// Get output as much as possible
	go func() {
		for {
			line, _, err := server.reader.ReadLine()
			if err != nil {
				log.Println("Error reading stdout from mindustry:", err)
			}
			log.Println("Mindustry server:", string(line))
			server.outputBuffer = append(server.outputBuffer, line...)
			server.outputBuffer = append(server.outputBuffer, '\n')
			server.outputChanged = true
		}
	}()

	return nil
}

func (server *MindustryServer) SendCommand(command string) (err error) {
	// clear output buffer
	server.outputBuffer = make([]byte, 0)
	_, err = server.inPipe.Write([]byte(command + "\n"))
	return err
}

// something like \033[101m    , use this to delete all the color code to make it plain text
var colorCodeReplace = regexp.MustCompile("\033" + regexp.QuoteMeta("[") + "[0-9]+m")

func (server *MindustryServer) GetOutput() (output []byte) {
	server.outputChanged = false
	// replace color code in terminal
	out := colorCodeReplace.ReplaceAll(server.outputBuffer, []byte(""))
	return out
}

func (server *MindustryServer) GetRawOutput() (output []byte) {
	server.outputChanged = false
	return server.outputBuffer
}

func (server MindustryServer) IsOutputUpdated() bool {
	return server.outputChanged
}

func (server *MindustryServer) Shutdown() (err error) {
	err = server.SendCommand("stop")
	if err != nil {
		return err
	}
	err = server.SendCommand("exit")
	if err != nil {
		return err
	}
	return nil
}
