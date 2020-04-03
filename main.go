package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	obsws "github.com/christopher-dG/go-obs-websocket"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Invalid number of args, specify command")
	}

	// Connect a client.
	c := obsws.Client{Host: "localhost", Port: 4444}
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer c.Disconnect()

	// Set the amount of time we can wait for a response.
	obsws.SetReceiveTimeout(time.Second * 2)

	var resp interface{}
	var err error

	cmds := []struct {
		regex string
		run   func(g []string)
	}{
		{"/scene=(.*)", func(g []string) {
			resp, err = obsws.NewSetCurrentSceneRequest(g[1]).SendReceive(c)
		}},
		{"/vol=([^,]+),([0-9.]+)", func(g []string) {
			level, err2 := strconv.ParseFloat(g[2], 64)
			if err2 != nil {
				err = err2
				return
			}
			resp, err = obsws.NewSetVolumeRequest(g[1], level).SendReceive(c)
		}},
	}

	for _, cmd := range cmds {
		if match := regexp.MustCompile(cmd.regex).FindStringSubmatch(strings.Join(os.Args[1:], " ")); match != nil {
			cmd.run(match)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	cnt, _ := json.MarshalIndent(resp, "", "  ")
	log.Println("set text:", string(cnt))
}
