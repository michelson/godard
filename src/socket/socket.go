package socket

import (
  "net"
  "log"
  "os"
  "os/signal"
  "syscall"
)

const Timeout = 60 // Used for client commands
const MaxAttempts = 5

type Socket struct{
    Listener net.Listener
}


func NewSocket() (*Socket, error) {
  // Create the socket to listen on:
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

func ClientCommand(base_dir string, name string , command string) (string, error){
  

  /*
    def client_command(base_dir, name, command)
      res = nil
      MAX_ATTEMPTS.times do |current_attempt|
        begin
          client(base_dir, name) do |socket|
            Timeout.timeout(TIMEOUT) do
              socket.puts command
              res = Marshal.load(socket.read)
            end
          end
          break
        rescue EOFError, Timeout::Error
          if current_attempt == MAX_ATTEMPTS - 1
            abort("Socket Timeout: Server may not be responding")
          end
          puts "Retry #{current_attempt + 1} of #{MAX_ATTEMPTS}"
        end
      end
      res
    end
  */
    res := "aa"

    return res, nil
}

/*
    def client(base_dir, name, &block)
      UNIXSocket.open(socket_path(base_dir, name), &block)
    end

    def client_command(base_dir, name, command)
      res = nil
      MAX_ATTEMPTS.times do |current_attempt|
        begin
          client(base_dir, name) do |socket|
            Timeout.timeout(TIMEOUT) do
              socket.puts command
              res = Marshal.load(socket.read)
            end
          end
          break
        rescue EOFError, Timeout::Error
          if current_attempt == MAX_ATTEMPTS - 1
            abort("Socket Timeout: Server may not be responding")
          end
          puts "Retry #{current_attempt + 1} of #{MAX_ATTEMPTS}"
        end
      end
      res
    end
*/