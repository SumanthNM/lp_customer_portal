/**
 * Functionalities
 * Maintain Active connections
 * callbacks for events
 * 	- on connection
 *	- on disconnection
 * Provide Read function - function to get bytes
**/
package websocket

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chassis/openlog"
	"github.com/gorilla/websocket"
)

type ReaderCallBack func(*context.Context) []byte // gets only

type Client struct {
	Pool     *Pool
	Conn     *websocket.Conn
	Context  *context.Context
	GetData  ReaderCallBack
	Interval int // interval in seconds.
	Ctx      *context.Context
}

func (c *Client) StartWriteLoop() {
	// send data
	openlog.Debug("Writing to client [" + c.Conn.RemoteAddr().String() + "]")
	data := c.GetData(c.Ctx)
	err := c.Conn.WriteMessage(1, data)
	if err != nil {
		openlog.Error("Error occured while writing [" + err.Error() + "]")
		c.Pool.Unregister <- c
		c.Conn.Close()
		return
	}
	// sleep for interval time
	time.Sleep(time.Duration(c.Interval) * time.Second)
	// time.Sleep(time.Duration(c.Interval * 1000))
	c.StartWriteLoop()
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type Pool struct {
	Register       chan *Client
	Unregister     chan *Client
	Clients        map[*Client]bool
	MaxConnections int64
}

func NewPool(maxConntections int64) *Pool {
	return &Pool{
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Clients:        make(map[*Client]bool),
		MaxConnections: maxConntections,
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			openlog.Debug("Connected successfully ")
			go client.StartWriteLoop() // start writing the data to client
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			openlog.Debug("client disconnected successfully ")
			break
		}
	}
}

var _upgrader *websocket.Upgrader

// takes create check origin function and initializes the upgrader.
func CreateUpgrader(CheckOrigin func(*http.Request) bool) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     CheckOrigin,
	}
	_upgrader = &upgrader
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := _upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}
