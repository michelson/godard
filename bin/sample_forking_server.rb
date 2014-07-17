#! /usr/bin/env ruby

# This is a modified version found at http://tomayko.com/writings/unicorn-is-unix
# It is modified to trigger various states like increase memory consumption so that
# I could write watches for them.

# Instructions for running the test
#
# (1) Edit the example config and fix the path to this file. Around line 16.
# (2) Load up the config and run the bluepill daemon
# (3) Run watch -n0.2 'sudo ruby bin/bluepill status 2>/dev/null; echo; ps ajxu | egrep "(CPU|forking|bluepill|sleep|ruby)" | grep -v grep | sort'
# (4) After verifying that the "sleep" workers are properly being restarted, telnet to localhost 4242 and say something. You should get it echoed back and the worker which answered your request should now be over the allowed memory limit
# (5) Observe the worker being killed in the watch you started in step 3.

require 'socket'

port = ARGV[0].to_i
port = 4242 if port == 0

acceptor = Socket.new(Socket::AF_INET, Socket::SOCK_STREAM, 0)
address = Socket.pack_sockaddr_in(port, '127.0.0.1')
acceptor.bind(address)
acceptor.listen(10)

children = []
trap('EXIT') { acceptor.close; children.each {|c| Process.kill('QUIT', c)} }


3.times do
  children << fork do
    trap('QUIT') {$0 = "forking_server| QUIT received shutting down gracefully..."; sleep 5; exit}
    trap('INT') {$0 = "forking_server| INT received shutting down UN-gracefully..."; sleep 3; exit}

    puts "child #$$ accepting on shared socket (localhost:#{port})"
    loop {
      socket, addr = acceptor.accept
      socket.write "child #$$ echo> "
      socket.flush
      message = socket.gets
      socket.write message
      socket.close
      puts "child #$$ echo'd: '#{message.strip}'"

      # cause a spike in mem usage
      temp = "*" * (100 * 1024)
    }
    exit
  end
end

trap('INT') { puts "\nbailing" ; exit }

Process.waitall
