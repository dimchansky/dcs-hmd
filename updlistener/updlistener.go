package updlistener

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

func New(port int) (*UPDListener, error) {
	pc, err := net.ListenPacket("udp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen udp packet: %w", err)
	}

	closeCh := make(chan struct{})
	l := &UPDListener{pc: pc, closeCh: closeCh}
	go l.listen()

	return l, nil
}

type UPDListener struct {
	closeOnce sync.Once
	closeErr  error
	closeCh   chan struct{}

	pc net.PacketConn
}

func (l *UPDListener) Close() error {
	l.closeOnce.Do(func() {
		l.closeErr = l.pc.Close()
		<-l.closeCh // wait until listener is stopped
	})

	return l.closeErr
}

func (l *UPDListener) listen() {
	closeCh := l.closeCh
	defer close(closeCh)

	pc := l.pc
	message := make([]byte, 8096)
	for {
		n, _, err := pc.ReadFrom(message[:])
		if err != nil {
			return
		}

		fmt.Println(string(message[:n]))
	}
}
