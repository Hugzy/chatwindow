package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func main() {
	log.SetFlags(0)

	server := MakeServer()

	fmt.Println("starting server on port 8000")
	//fn := http.HandlerFunc(server.Handle)
	http.HandleFunc("/connect/{group...}", server.Handle)
	err := func() error {
		server := &http.Server{Addr: "localhost:8000"}
		return server.ListenAndServe()
	}()

	if err != nil {
		log.Fatal(err)
	}

}

type Connection struct {
	conn      *websocket.Conn
	connected bool
}

type Server struct {
	conns map[string][]*Connection
}

func MakeServer() *Server {
	srv := &Server{
		conns: make(map[string][]*Connection),
	}
	return srv
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url is", r.URL)

	group := r.PathValue("group")

	if group == "" {
		group = "default"
	}

	fmt.Println("Connect to group", group)

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})

	if err != nil {
		fmt.Println("cannot accept message from websocket", err)
	}

	conn := &Connection{
		conn:      c,
		connected: true,
	}

	s.conns[group] = append(s.conns[group], conn)

	s.Serve(group, c)
}

func (s *Server) Broadcast(group string, msg []byte) {
	conns := s.conns[group]
	for _, conn := range conns {
		if err := conn.conn.Write(context.TODO(), websocket.MessageText, msg); err != nil {
			fmt.Println("error writing message", err)
		}
	}
}

func (s *Server) Serve(group string, ws *websocket.Conn) error {
	buf := make([]byte, 1024)
	for {
		_, reader, err := ws.Reader(context.TODO())
		if err != nil {
			//fmt.Println("could not create reader", err)
			continue
		}

		n, err := reader.Read(buf)
		if err != nil {
			if err == io.ErrClosedPipe {
				fmt.Println("read error", err)
				continue
			}
		}
		msg := buf[:n]
		s.Broadcast(group, msg)
	}
}
