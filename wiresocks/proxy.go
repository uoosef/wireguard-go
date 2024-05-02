package wiresocks

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/netip"
	"time"

	"github.com/bepass-org/warp-plus/proxy/pkg/mixed"
	"github.com/bepass-org/warp-plus/proxy/pkg/statute"
	"github.com/bepass-org/warp-plus/wireguard/device"
	"github.com/bepass-org/warp-plus/wireguard/tun/netstack"
	"github.com/things-go/go-socks5/bufferpool"
)

// VirtualTun stores a reference to netstack network and DNS configuration
type VirtualTun struct {
	Tnet   *netstack.Net
	Logger *slog.Logger
	Dev    *device.Device
	Ctx    context.Context
	pool   bufferpool.BufPool
}

// StartProxy spawns a socks5 server.
func (vt *VirtualTun) StartProxy(bindAddress netip.AddrPort) (netip.AddrPort, error) {
	ln, err := net.Listen("tcp", bindAddress.String())
	if err != nil {
		return netip.AddrPort{}, err // Return error if binding was unsuccessful
	}

	proxy := mixed.NewProxy(
		mixed.WithListener(ln),
		mixed.WithLogger(vt.Logger),
		mixed.WithContext(vt.Ctx),
		mixed.WithUserHandler(func(request *statute.ProxyRequest) error {
			return vt.generalHandler(request)
		}),
	)
	go func() {
		_ = proxy.ListenAndServe()
	}()
	go func() {
		<-vt.Ctx.Done()
		vt.Stop()
	}()

	return ln.Addr().(*net.TCPAddr).AddrPort(), nil
}

func (vt *VirtualTun) generalHandler(req *statute.ProxyRequest) error {
	vt.Logger.Info("handling connection", "protocol", req.Network, "destination", req.Destination)
	conn, err := vt.Tnet.Dial(req.Network, req.Destination)
	if err != nil {
		return err
	}
	// Close the connections when this function exits
	defer conn.Close()
	defer req.Conn.Close()
	// Channel to notify when copy operation is done
	done := make(chan error, 1)
	// Copy data from req.Conn to conn
	go func() {
		req.Conn.SetReadDeadline(time.Now().Add(15 * time.Second))
		buf1 := vt.pool.Get()
		defer vt.pool.Put(buf1)
		_, err := copyConnTimeout(conn, req.Conn, buf1[:cap(buf1)], 15*time.Second)
		done <- err
	}()
	// Copy data from conn to req.Conn
	go func() {
		conn.SetReadDeadline(time.Now().Add(15 * time.Second))
		buf2 := vt.pool.Get()
		defer vt.pool.Put(buf2)
		_, err := copyConnTimeout(req.Conn, conn, buf2[:cap(buf2)], 15*time.Second)
		done <- err
	}()
	// Wait for one of the copy operations to finish
	err = <-done
	if err != nil {
		vt.Logger.Warn(err.Error())
	}

	// Close connections and wait for the other copy operation to finish
	<-done
	return nil
}

func (vt *VirtualTun) Stop() {
	if vt.Dev != nil {
		if err := vt.Dev.Down(); err != nil {
			vt.Logger.Warn(err.Error())
		}
	}
}

var errInvalidWrite = errors.New("invalid write result")

func copyConnTimeout(dst net.Conn, src net.Conn, buf []byte, timeout time.Duration) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in CopyBuffer")
	}

	for {
		if err := src.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return 0, err
		}

		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
