package updlistener

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

type MessageHandler interface {
	HandleMessage(msg []byte)
}

type MessageHandlerFunc func(msg []byte)

func (f MessageHandlerFunc) HandleMessage(msg []byte) {
	f(msg)
}

func New(port int, msgHandler MessageHandler) (*UPDListener, error) {
	pc, err := net.ListenPacket("udp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen udp packet: %w", err)
	}

	closeCh := make(chan struct{})
	l := &UPDListener{closeCh: closeCh, pc: pc, msgHandler: msgHandler}

	go l.listen()

	return l, nil
}

type UPDListener struct {
	closeOnce sync.Once
	closeErr  error
	closeCh   chan struct{}

	pc         net.PacketConn
	msgHandler MessageHandler
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

	const bufSize = 8096
	message := make([]byte, bufSize)
	pc := l.pc
	msgHandler := l.msgHandler

	for {
		n, _, err := pc.ReadFrom(message)
		if err != nil {
			return
		}

		msgHandler.HandleMessage(message[:n])
	}
}
