require 'socket'
sock = UDPSocket.new
sock.bind('0.0.0.0', 9121)
sock.send([1, 2, 0x10, 0, 0, 0, 0, 0, 0, 0, 0, 0].pack('C*'), 0, '8.8.8.8', 53)
puts sock.recvfrom(20)
