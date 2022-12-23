package mindustryserver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"time"
)

type MindustryServer struct {
	cmd           *exec.Cmd
	running       bool
	inPipe        io.WriteCloser
	outPipe       io.ReadCloser
	reader        *bufio.Reader
	outputChannel chan []byte
}

func NewMindustryServer() *MindustryServer {
	server := new(MindustryServer)

	server.cmd = exec.Command("java", "-jar", "./mindustry-server/server.jar")
	server.running = false
	server.inPipe, _ = server.cmd.StdinPipe()
	server.outPipe, _ = server.cmd.StdoutPipe()
	server.outputChannel = make(chan []byte)

	return server
}

// something like \033[101m    , use this to delete all the color code to make it plain text
var colorCodeReplace = regexp.MustCompile("\033" + regexp.QuoteMeta("[") + "[0-9]+m")

func (server *MindustryServer) Start() (err error) {
	if server.running {
		return
	}
	err = server.cmd.Start()
	if err != nil {
		return err
	}

	server.reader = bufio.NewReader(server.outPipe)
	server.running = true

	// Get output as much as possible
	go func() {
		for server.running {
			line, _, err := server.reader.ReadLine()
			if err != nil {
				// this happens when server was killed
				// should have a better way to handle this
				if err == io.EOF {
					time.Sleep(time.Millisecond * 500)
					continue
				}
				log.Println("Error reading stdout from mindustry:", err)
			}

			// also print output to stdout
			fmt.Println(string(line))
			line = colorCodeReplace.ReplaceAll(line, []byte(""))
			line = append(line, byte('\n'))
			server.outputChannel <- line
		}
	}()

	return nil
}

func (server *MindustryServer) SendCommand(command string) (err error) {
	_, err = server.inPipe.Write([]byte(command + "\n"))
	return err
}

func (server MindustryServer) GetOutputChannel() chan []byte {
	return server.outputChannel
}

func (server *MindustryServer) Exit() (err error) {
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

func (server *MindustryServer) Kill() (err error) {
	server.running = false
	err = server.cmd.Process.Kill()
	return err
}
