package mindustryserver

import (
	"bufio"
	"io"
	"os/exec"
)

type MindustryServer struct {
	cmd          *exec.Cmd
	started      bool
	inPipe       io.WriteCloser
	outPipe      io.ReadCloser
	scanner      *bufio.Scanner
	outputBuffer []byte
}

func NewMindustryServer() MindustryServer {
	server := *new(MindustryServer)

	server.cmd = exec.Command("java", "-jar", "./mindustry-server/server.jar")
	server.started = false
	server.inPipe, _ = server.cmd.StdinPipe()
	server.outPipe, _ = server.cmd.StdoutPipe()
	server.outputBuffer = make([]byte, 0)
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

	server.scanner = bufio.NewScanner(server.outPipe)
	server.started = true

	// Get output as much as possible
	go func() {
		for server.scanner.Scan() {
			server.outputBuffer = append(server.outputBuffer, server.scanner.Bytes()...)
			server.outputBuffer = append(server.outputBuffer, '\n')
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

func (server MindustryServer) GetOutput() (output []byte) {
	return server.outputBuffer
}

func (server MindustryServer) Shutdown() (err error) {
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
