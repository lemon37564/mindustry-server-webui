package mindustryserver

import (
	"bufio"
	"io"
	"os/exec"
)

type MindustryServer struct {
	cmd     *exec.Cmd
	inPipe  io.WriteCloser
	outPipe io.ReadCloser
	reader  *bufio.Reader
}

func NewMindustryServer() MindustryServer {
	server := *new(MindustryServer)

	server.cmd = exec.Command("java", "-jar", "./mindustry-server/server.jar")

	server.inPipe, _ = server.cmd.StdinPipe()
	server.outPipe, _ = server.cmd.StdoutPipe()
	return server
}

func (server *MindustryServer) Start() {
	err := server.cmd.Start()
	if err != nil {
		panic(err)
	}

	server.reader = bufio.NewReader(server.outPipe)
}

func (server MindustryServer) SendCommand(command string) (err error) {
	_, err = server.inPipe.Write([]byte(command + "\n"))
	return err
}

func (server MindustryServer) GetOutput() (output []byte, err error) {
	output, _, err = server.reader.ReadLine()
	return output, err
}

func (server MindustryServer) Shutdown() (err error) {
	return server.cmd.Process.Kill()
}
