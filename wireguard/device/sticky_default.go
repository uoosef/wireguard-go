//go:build !linux

package device

import (
	"github.com/bepass-org/wireguard-go/wireguard/conn"
	"github.com/bepass-org/wireguard-go/wireguard/rwcancel"
)

func (device *Device) startRouteListener(bind conn.Bind) (*rwcancel.RWCancel, error) {
	return nil, nil
}
