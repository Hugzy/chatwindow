package main

import (
	"context"
	"errors"
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
  fn := http.HandlerFunc(server.Handle)
  err := http.ListenAndServe("localhost:8000", fn)
  
  if err!= nil {
    log.Fatal(err)
  }

}

type Server struct {
  conns map[*websocket.Conn]bool
}

func MakeServer() *Server {
  srv := &Server {
    conns: make(map[*websocket.Conn]bool),
  }
  return srv
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
  fmt.Println("new request")
  
  c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})

  if err != nil{
    fmt.Println("cannot accept message from websocket", err)
  }

  s.conns[c] = true

  s.Serve(c)
}

func (s *Server) Serve(ws *websocket.Conn) error {
  buf := make([]byte, 1024)
  for {
    t, reader, err := ws.Reader(context.TODO())
    if err != nil {
       fmt.Println("could not create reader", err)
    }

    n, err := reader.Read(buf)
    if err != nil {
      if err == io.ErrClosedPipe {
        fmt.Println("read error", err)
        continue
      }
    }
  msg := buf[:n]
  fmt.Println(string(msg))
  ws.Write(context.TODO(), t, []byte("thank you for the message"))
  }
}


func run() error {
  /*
  errc := make(chan error, 1)
  go func() {
    errc <- nil
  }()

  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, os.Interrupt)
  select {
  case err := <-errc:
    log.Printf("failed to serve: %v", err)
  case sig := <-sigs:
    log.Printf("terminating: %v", sig)
  }

  ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
  defer cancel()

  return s.Shutdown(ctx)
  */
  return errors.New("wtf?")
}
