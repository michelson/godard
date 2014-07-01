package socket

import (
  "net"
  "log"
  "os"
  "os/signal"
  "syscall"
)

type Socket struct{
    //TIMEOUT = 60 # Used for client commands
    //MAX_ATTEMPTS = 5
    Listener net.Listener
}


func NewSocket() (*Socket, error) {
  // Create the socket to listen on:
    os.Remove("/tmp/godard.sock") // just in case
    l, err := net.Listen("unix", "/tmp/godard.sock")
    c := &Socket{}
    if err != nil {
      log.Fatal(err)
      return c, err
    }

    c.Listener = l


    return c, err
}


func (c *Socket) Run() {

  
  sigc := make(chan os.Signal, 1)
  signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
  go func(cc chan os.Signal) {
      // Wait for a SIGINT or SIGKILL:
      sig := <-cc
      log.Printf("Caught signal %s: shutting down.", sig)
      // Stop listening (and unlink the socket if unix type):
      c.Listener.Close()
      os.Remove("/tmp/godard.sock")
      // And we're done:
      os.Exit(0)
  }(sigc)
  

  for {
    log.Println("ACTION!")
    fd, err := c.Listener.Accept()
    if err != nil {
      log.Println("accept error", err)
    }
    go c.EchoServer(fd)
  }

}


func (c *Socket) EchoServer(conn net.Conn) {
  for {
      buf := make([]byte, 512)
      nr, err := conn.Read(buf)
      if err != nil {
          return
      }

      data := buf[0:nr]
      println("Server got:", string(data))
      _, err = conn.Write(data)
      if err != nil {
          panic(err)
      }
  }
}
