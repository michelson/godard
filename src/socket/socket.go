package socket

import (
	"log"
	"net"
	"io"
	"path"
)

//http://stackoverflow.com/questions/2886719/unix-sockets-in-go
const Timeout = 60 // Used for client commands
const MaxAttempts = 5

var messages chan string = make(chan string)

type Socket struct {
	Listener        net.Listener
	ListenerChannel chan string
	Path            string
}

func NewSocket(base_dir string, name string) (*Socket, error) {
	// Create the socket to listen on:
	log.Println("SOCKET PATH ", base_dir)
	c := &Socket{}
	c.Path = SocketPath(base_dir, name)
	l, err := net.Listen("unix", c.Path)
	if err != nil {
		log.Fatal(err)
		return c, err
	}

	c.Listener = l
	c.ListenerChannel = messages

	return c, err
}

func (c *Socket) Run() {
	/*
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
	*/

	for {
		//log.Println("ACTION!")
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
		//println("Server got:", string(data))
		/*
		 */
		go func(data string) {
			c.ListenerChannel <- string(data)
		}(string(data))

		_, err = conn.Write(data)
		if err != nil {
			log.Println("SOCK ERR:", err)
			//panic(err)
		}
	}
}

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		messages <- string(buf[0:n])
		//println("Client got:", string(buf[0:n]))
	}
}

func ClientCommand(base_dir string, name string, command string) (string, error) {

	log.Println("PREPARE COMMAND TO SEND:", base_dir, name, command)

	c, err := net.Dial("unix", SocketPath(base_dir, name))

	if err != nil {
		panic(err)
	}

	defer c.Close()

	go reader(c)

	// set retries attempts and everything

	_, err2 := c.Write([]byte(command))
	if err2 != nil {
		println(err2)
		//break
	}

	msg := <-messages
	//println("RES!!", msg)

	res := msg

	return res, nil
}

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

func SocketPath(base_dir string, name string) string {
	s := path.Join(base_dir, "sock", name+".sock")
	return s
}
