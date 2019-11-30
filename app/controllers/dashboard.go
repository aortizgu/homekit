package controllers

import (
	"homekit/app/routes"
	"log"

	"homekit/app/msgbroker"

	"github.com/revel/revel"
)

type Dashboard struct {
	Application
}

func (c Dashboard) checkUser() revel.Result {
	log.Println("checkuser")
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Dashboard) Index() revel.Result {
	c.Log.Info("Fetching index")
	return c.Render()
}

func (c Dashboard) Live(user string, ws revel.ServerWebSocket) revel.Result {
	// Make sure the websocket is valid.
	if ws == nil {
		return nil
	}

	// Join the room.
	subscription := msgbroker.Subscribe()
	defer subscription.Cancel()

	msgbroker.Join(user)
	defer msgbroker.Leave(user)

	// Send down the archive.
	for _, event := range subscription.Archive {
		if ws.MessageSendJSON(&event) != nil {
			// They disconnected
			return nil
		}
	}

	// In order to select between websocket messages and subscription events, we
	// need to stuff websocket events into a channel.
	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := ws.MessageReceiveJSON(&msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	// Now listen for new events from either the websocket or the msgbroker.
	for {
		select {
		case event := <-subscription.New:
			if ws.MessageSendJSON(&event) != nil {
				// They disconnected.
				return nil
			}
		case msg, ok := <-newMessages:
			// If the channel is closed, they disconnected.
			if !ok {
				return nil
			}

			// @aortiz: we don't want say nothing in this app
			// Otherwise, say something.
			// msgbroker.Say(user, msg)
			log.Printf("ws client says[%s]\n", msg)
		}
	}

	return nil
}
