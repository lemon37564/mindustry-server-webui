package mindustryserver

import (
	"bufio"
	"io"
	"os/exec"
)

type MindustryServer struct {
	cmd     *exec.Cmd
	started bool
	inPipe  io.WriteCloser
	outPipe io.ReadCloser
	reader  *bufio.Reader
}

func NewMindustryServer() MindustryServer {
	server := *new(MindustryServer)

	server.cmd = exec.Command("java", "-jar", "./mindustry-server/server.jar")
	server.started = false
	server.inPipe, _ = server.cmd.StdinPipe()
	server.outPipe, _ = server.cmd.StdoutPipe()
	return server
}

func (server *MindustryServer) Start() {
	if server.started {
		return
	}
	err := server.cmd.Start()
	if err != nil {
		panic(err)
	}

	server.reader = bufio.NewReader(server.outPipe)
	server.started = true
}

func (server MindustryServer) SendCommand(command string) (err error) {
	_, err = server.inPipe.Write([]byte(command + "\n"))
	return err
}

func (server MindustryServer) GetOutput() (output []byte, err error) {
	output = make([]byte, 4096)
	_, err = server.reader.Read(output)
	return output, err
}

func (server MindustryServer) Shutdown() (err error) {
	return server.cmd.Process.Kill()
}
