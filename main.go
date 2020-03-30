package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	obsws "github.com/christopher-dG/go-obs-websocket"
)

func main() {
	// Connect a client.
	c := obsws.Client{Host: "localhost", Port: 4444}
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer c.Disconnect()

	req := obsws.NewGetTextFreetype2PropertiesRequest("IP")
	resp, err := req.SendReceive(c)
	if err != nil {
		log.Fatal(err)
	}

	cnt, _ := json.MarshalIndent(resp, "", "  ")
	log.Println("streaming:", string(cnt))

	// Set the amount of time we can wait for a response.
	obsws.SetReceiveTimeout(time.Second * 2)

	go func() {
		for i := 0; ; i++ {
			time.Sleep(5 * time.Second)

			req := obsws.NewSetTextFreetype2PropertiesRequest(
				"IP",
				0xeeeeeeff, 0xeeeeeeff, 0, false, nil, "", 0, 0, "", false, false, false,
				fmt.Sprintf("This is the text %d", i),
				"", false,
			)
			resp, err := req.SendReceive(c)
			if err != nil {
				log.Fatal(err)
			}

			cnt, _ := json.MarshalIndent(resp, "", "  ")
			log.Println("set text:", string(cnt))
		}
	}()

	// Respond to events by registering handlers.
	c.MustAddEventHandler("SwitchScenes", func(e obsws.Event) {
		// Make sure to assert the actual event type.
		log.Println("new scene:", e.(obsws.SwitchScenesEvent).SceneName)
	})
	c.MustAddEventHandler("TransitionBegin", func(e obsws.Event) {
		// Make sure to assert the actual event type.
		ev := e.(obsws.TransitionBeginEvent)
		log.Println("transition begin:", ev.Name, ev.Duration, ev.FromScene, ev.ToScene)
	})

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	s := <-signals
	switch s {
	case syscall.SIGTERM, syscall.SIGINT:
		log.Println("Received Ctrl+C")
	default:
		log.Println("Received signal", s)
	}
}
