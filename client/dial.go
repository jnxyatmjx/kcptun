package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"

	kcp "github.com/jnxyatmjx/kcp-go"
	"github.com/jnxyatmjx/kcptun/generic"
	"github.com/pkg/errors"
	"github.com/xtaci/tcpraw"
)

func dial(config *Config, block kcp.BlockCrypt) (*kcp.UDPSession, error) {
	mp, err := generic.ParseMultiPort(config.RemoteAddr)
	if err != nil {
		return nil, err
	}

	var randport uint64
	err = binary.Read(rand.Reader, binary.LittleEndian, &randport)
	if err != nil {
		return nil, err
	}

	remoteAddr := fmt.Sprintf("%v:%v", mp.Host, uint64(mp.MinPort)+randport%uint64(mp.MaxPort-mp.MinPort+1))

	if config.TCP {
		conn, err := tcpraw.Dial("tcp", remoteAddr)
		if err != nil {
			return nil, errors.Wrap(err, "tcpraw.Dial()")
		}
		return kcp.NewConn(remoteAddr, block, config.DataShard, config.ParityShard, conn)
	}
	return kcp.DialWithOptions(remoteAddr, block, config.DataShard, config.ParityShard)

}
