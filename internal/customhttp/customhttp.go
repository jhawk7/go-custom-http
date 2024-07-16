package customhttp

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CustomConn interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type Chttp struct {
	conn   CustomConn
	routes map[string]string
}

func InitHttp(c CustomConn, routes map[string]string) *Chttp {
	return &Chttp{
		conn:   c,
		routes: routes,
	}
}

func (c *Chttp) HandleRequest() {
	buf := make([]byte, 1024)
	end, rErr := c.conn.Read(buf)
	if rErr != nil {
		err := fmt.Errorf("failed to read request from conn; %v", rErr)
		log.Error(err)
		fmt.Fprintf(c.conn, "HTTP/1.1 500 Internal Service Error")
		return
	}

	request := string(buf[:end])
	log.Infof("request received; length %v", end)
	parts := strings.SplitN(request, " ", 3)
	method := parts[0]
	path := parts[1]
	var data string
	if len(parts) == 3 {
		data = parts[2]
	}

	switch method {
	case "GET":
		c.getHandler(path)
	case "POST":
		c.postHandler(path, data)
	case "DELETE":
		c.deleteHandler(path, data)
	default:
		fmt.Fprintf(c.conn, "HTTP/1.1 405 Method Not Allowed\r\nContent-Length:24\r\n\r\nInvalid request method.")
	}

	defer func() {
		c.conn.Close()
		log.Info("connection closed")
	}()
}

func (c *Chttp) getHandler(path string) {
	if path != c.routes["GET"] {
		msg := "Invalid url."
		fmt.Fprintf(c.conn, "HTTP/1.1 404 Not Found\r\nContent-Length: %v\r\n\r\n%v", len(msg), msg)
		return
	}

	log.Info("GET request received")
	msg := "<h1>Hello!</h1>"
	fmt.Fprintf(c.conn, "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n Content-Length: %v\r\n\r\n%v", len(msg), msg)
}

func (c *Chttp) postHandler(path string, data string) {
	if path != c.routes["POST"] {
		msg := "Invalid url."
		fmt.Fprintf(c.conn, "HTTP/1.1 404 Not Found\r\nContent-Length: %v\r\n\r\n%v", len(msg), msg)
		return
	}

	log.Infof("POST request received with data %v", data)
	fmt.Fprintf(c.conn, "HTTP/1.1 201 Created")
}

func (c *Chttp) deleteHandler(path string, data string) {
	if path != c.routes["DELETE"] {
		msg := "Invalid url."
		fmt.Fprintf(c.conn, "HTTP/1.1 404 Not Found\r\nContent-Length: %v\r\n\r\n%v", len(msg), msg)
	}

	log.Infof("DELETE request received with data %v", data)
	fmt.Fprintf(c.conn, "HTTP/1.1 204 No Content")
}
