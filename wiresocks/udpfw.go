package wiresocks

import (
	"context"
	"net"
	"sync"
)

func NewVtunUDPForwarder(ctx context.Context, localBind, dest string, vtun *VirtualTun, mtu int) error {
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
