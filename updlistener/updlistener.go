package updlistener

import (
	"bufio"
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
	address := ":" + strconv.Itoa(port)

	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve udp address '%s': %w", address, err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen udp packet: %w", err)
	}

	closeCh := make(chan struct{})
	l := &UPDListener{closeCh: closeCh, conn: conn, msgHandler: msgHandler}

	go l.listen()

	return l, nil
}

type UPDListener struct {
	closeOnce sync.Once
	closeErr  error
	closeCh   chan struct{}

	conn       *net.UDPConn
	msgHandler MessageHandler
}

func (l *UPDListener) Close() error {
	l.closeOnce.Do(func() {
		l.closeErr = l.conn.Close()
		<-l.closeCh // wait until listener is stopped
	})

	return l.closeErr
}

func (l *UPDListener) listen() {
	closeCh := l.closeCh
	defer close(closeCh)

	const bufSize = 64 * 1024

	reader := bufio.NewReaderSize(l.conn, bufSize)
	msgHandler := l.msgHandler

	readLines(reader, msgHandler)
}

func readLines(reader *bufio.Reader, msgHandler MessageHandler) {
	nextIsContinuation := false

	var (
		message []byte
		err     error
	)

	for {
		isContinuation := nextIsContinuation

		message, nextIsContinuation, err = reader.ReadLine()
		if err != nil {
			return
		}

		if isContinuation || nextIsContinuation { // skip lines that do not fit in the buffer
			continue
		}

		msgHandler.HandleMessage(message)
	}
}
