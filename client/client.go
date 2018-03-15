package client

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/orvice/kit/log"
	"github.com/orvice/monitor-server/mod"
)

type Client struct {
	wsAddr string
	NodeID int32
	wsConn *websocket.Conn
	ctx    context.Context
	cancel func()

	done chan int32

	messageRecv chan mod.Packet
	MessageSend chan []byte
	logger      log.Logger
}

func NewClient(nodeID int32, wsAddr string, msgRecv chan mod.Packet,
	l log.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		wsAddr:      wsAddr,
		NodeID:      nodeID,
		ctx:         ctx,
		cancel:      cancel,
		messageRecv: msgRecv,
		MessageSend: make(chan []byte),
		logger:      l,
	}
	return c
}

func (c *Client) Exit() {
	c.cancel()
	c.wsConn.Close()
	c.done <- c.NodeID
}

func (c *Client) Reconnect() {
	for {
		err := c.Connect()
		if err == nil {
			break
		}
	}
	c.logger.Infof("reconnect success %d", c.NodeID)
	go c.Read()
	go c.Write()
}

func (c *Client) Connect() error {
	var err error
	c.logger.Infof("connecting to %s", c.wsAddr)
	c.wsConn, _, err = websocket.DefaultDialer.Dial(c.wsAddr, nil)
	if err != nil {
		c.logger.Errorf("dial %s : error %v", c.wsAddr, err)
		time.Sleep(time.Second * 2)
		return err
	}
	return nil
}

func (c *Client) Run() {
	err := c.Connect()
	if err != nil {
		c.Reconnect()
		return
	}
	go c.Read()
	go c.Write()
}

func (c *Client) Read() {

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, message, err := c.wsConn.ReadMessage()
			if err != nil {
				c.logger.Error("read:", err)
				c.Reconnect()
				return
			}
			c.messageRecv <- mod.Packet{
				NodeID:  c.NodeID,
				Message: message,
			}
		}
	}

}

func (c *Client) Write() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("ctx done...")
			err := c.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				c.logger.Error("write close:", err)
				return
			}
			return
		case t := <-ticker.C:
			err := c.wsConn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				c.logger.Error("write:", err)
				c.Reconnect()
				return
			}
		}
	}
}
