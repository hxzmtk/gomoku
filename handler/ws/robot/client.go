package robot

import "github.com/hxzmtk/gomoku/service/hub"

var (
	_ hub.IClient = &Client{}
)

type Client struct {
}

func (c *Client) ReadPump() {

}

func (c *Client) WritePump() {

}

func (c *Client) GetRoom() *hub.Room {
	return nil
}
func (c *Client) SetRoom(room *hub.Room) {
}

func (c *Client) GetID() string {
	return ""
}
func (c *Client) CloseChan() {

}
