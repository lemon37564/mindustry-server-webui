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

type mindustryBridge struct {
	cmd     *exec.Cmd
	running bool
	inPipe  io.WriteCloser
	outPipe io.ReadCloser
	reader  *bufio.Reader
	// a set of channel which handle bytes
	outputChannels map[chan []byte]struct{}
}

func newMindustryServer() *mindustryBridge {
	server := new(mindustryBridge)

	server.cmd = exec.Command("java", "-jar", "./mindustry-server/server.jar")
	server.running = false
	server.inPipe, _ = server.cmd.StdinPipe()
	server.outPipe, _ = server.cmd.StdoutPipe()
	server.outputChannels = make(map[chan []byte]struct{})

	return server
}

// something like \033[101m    , use this to delete all the color code to make it plain text
var colorCodeReplace = regexp.MustCompile("\033" + regexp.QuoteMeta("[") + "[0-9]+m")

func (server *mindustryBridge) start() (err error) {
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

			for ch := range server.outputChannels {
				ch <- line
			}
		}
	}()

	return nil
}

func (server *mindustryBridge) sendCommand(command string) (err error) {
	_, err = server.inPipe.Write([]byte(command + "\n"))
	return err
}

func (server *mindustryBridge) appendOutputChannel(ch chan []byte) {
	server.outputChannels[ch] = struct{}{}
}

func (server *mindustryBridge) RemoveOutputChannel(ch chan []byte) {
	delete(server.outputChannels, ch)
}

func (server *mindustryBridge) exit() (err error) {
	err = server.sendCommand("stop")
	if err != nil {
		return err
	}
	err = server.sendCommand("exit")
	if err != nil {
		return err
	}
	return nil
}

func (server *mindustryBridge) kill() (err error) {
	if server.cmd == nil {
		return nil
	}
	server.running = false
	err = server.cmd.Process.Kill()
	return err
}
