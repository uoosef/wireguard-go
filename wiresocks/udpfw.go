package wiresocks

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

type Socks5UDPForwarder struct {
	socks5Server string
	destAddr     *net.UDPAddr
	proxyUDPAddr *net.UDPAddr
	conn         *net.UDPConn
	listener     *net.UDPConn
	clientAddr   *net.UDPAddr
}

func NewVtunUDPForwarder(localBind, dest string, vtun *VirtualTun, mtu int, ctx context.Context) error {
	localAddr, err := net.ResolveUDPAddr("udp", localBind)
	if err != nil {
		return err
	}

	destAddr, err := net.ResolveUDPAddr("udp", dest)
	if err != nil {
		return err
	}

	listener, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return err
	}

	rconn, err := vtun.Tnet.DialUDP(nil, destAddr)
	if err != nil {
		return err
	}

	var clientAddr *net.UDPAddr
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		buffer := make([]byte, mtu)
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				n, cAddr, err := listener.ReadFrom(buffer)
				if err != nil {
					continue
				}

				clientAddr = cAddr.(*net.UDPAddr)

				rconn.WriteTo(buffer[:n], destAddr)
			}
		}
	}()
	go func() {
		buffer := make([]byte, mtu)
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				n, _, err := rconn.ReadFrom(buffer)
				if err != nil {
					continue
				}
				if clientAddr != nil {
					listener.WriteTo(buffer[:n], clientAddr)
				}
			}
		}
	}()
	go func() {
		wg.Wait()
		_ = listener.Close()
		_ = rconn.Close()
	}()
	return nil
}

func NewSocks5UDPForwarder(localBind, socks5Server, dest string) (*Socks5UDPForwarder, error) {
	localAddr, err := net.ResolveUDPAddr("udp", localBind)
	if err != nil {
		return nil, err
	}

	destAddr, err := net.ResolveUDPAddr("udp", dest)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}

	tcpConn, err := net.Dial("tcp", socks5Server)
	if err != nil {
		return nil, err
	}
	defer tcpConn.Close()

	if err := socks5Handshake(tcpConn); err != nil {
		return nil, err
	}

	proxyUDPAddr, err := requestUDPAssociate(tcpConn)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.DialUDP("udp", nil, proxyUDPAddr)
	if err != nil {
		return nil, err
	}

	return &Socks5UDPForwarder{
		socks5Server: socks5Server,
		destAddr:     destAddr,
		proxyUDPAddr: proxyUDPAddr,
		conn:         udpConn,
		listener:     listener,
	}, nil
}

func (f *Socks5UDPForwarder) Start() {
	go f.listenAndServe()
	go f.receiveFromProxy()
}

func socks5Handshake(conn net.Conn) error {
	// Send greeting
	_, err := conn.Write([]byte{0x05, 0x01, 0x00}) // SOCKS5, 1 authentication method, No authentication
	if err != nil {
		return err
	}

	// Receive server response
	resp := make([]byte, 2)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return err
	}

	if resp[0] != 0x05 || resp[1] != 0x00 {
		return fmt.Errorf("invalid SOCKS5 authentication response")
	}
	return nil
}

func (f *Socks5UDPForwarder) listenAndServe() {
	for {
		buffer := make([]byte, 4096)
		// Listen for incoming UDP packets
		n, clientAddr, err := f.listener.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from listener: %v\n", err)
			continue
		}

		// Store client address for response mapping
		f.clientAddr = clientAddr

		// Forward packet to destination via SOCKS5 proxy
		go f.forwardPacketToRemote(buffer[:n])
	}
}

func (f *Socks5UDPForwarder) forwardPacketToRemote(data []byte) {
	packet := make([]byte, 10+len(data))
	packet[0] = 0x00 // Reserved
	packet[1] = 0x00 // Reserved
	packet[2] = 0x00 // Fragment
	packet[3] = 0x01 // Address type (IPv4)
	copy(packet[4:8], f.destAddr.IP.To4())
	binary.BigEndian.PutUint16(packet[8:10], uint16(f.destAddr.Port))
	copy(packet[10:], data)

	_, err := f.conn.Write(packet)
	if err != nil {
		fmt.Printf("Error forwarding packet to remote: %v\n", err)
	}
}

func (f *Socks5UDPForwarder) receiveFromProxy() {
	for {
		buffer := make([]byte, 4096)
		n, err := f.conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from proxy connection: %v\n", err)
			continue
		}

		// Forward the packet to the original client
		f.listener.WriteToUDP(buffer[10:n], f.clientAddr)
	}
}

func requestUDPAssociate(conn net.Conn) (*net.UDPAddr, error) {
	// Send UDP associate request with local address and port set to zero
	req := []byte{0x05, 0x03, 0x00, 0x01, 0, 0, 0, 0, 0, 0} // Command: UDP Associate
	if _, err := conn.Write(req); err != nil {
		return nil, err
	}

	// Receive response
	resp := make([]byte, 10)
	if _, err := io.ReadFull(conn, resp); err != nil {
		return nil, err
	}

	if resp[1] != 0x00 {
		return nil, fmt.Errorf("UDP ASSOCIATE request failed")
	}

	// Parse the proxy UDP address
	bindIP := net.IP(resp[4:8])
	bindPort := binary.BigEndian.Uint16(resp[8:10])

	return &net.UDPAddr{IP: bindIP, Port: int(bindPort)}, nil
}
