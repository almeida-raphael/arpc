package channel

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"strconv"
	"strings"
)

// QUICStream is a channel stream implementation for QUIC
type QUICStream struct {
	stream quic.Stream
}

// Close closes a QUICStream
func (s *QUICStream) Close() error {
	return s.stream.Close()
}

// Write writes data on a QUICStream
func (s *QUICStream) Write(buf []byte) (int, error) {
	return s.stream.Write(buf)
}

// Read reads data from a QUICStream
func (s *QUICStream) Read(buf []byte) (int, error) {
	return s.stream.Read(buf)
}

// QUICSession is a session implementation for QUIC
type QUICSession struct {
	session quic.Session
}

// RemoteAddress returns address and port for a QUICSession
func (qs *QUICSession) RemoteAddress() (string, int, error) {
	// TODO: This function is only supporting IPV4
	addrParts := strings.Split(qs.session.RemoteAddr().String(), ":")
	if len(addrParts) != 2 {
		return "", 0, fmt.Errorf("invalid remote address")
	}

	portInt64, err := strconv.ParseInt(addrParts[1], 10, 32)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port number")
	}

	return addrParts[0], int(portInt64), nil
}

// AcceptStream receives a new stream from a client
func (qs *QUICSession) AcceptStream(ctx context.Context) (Stream, error) {
	stream, err := qs.session.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	return &QUICStream{stream: stream}, nil
}

// OpenStream creates a new stream with the server
func (qs *QUICSession) OpenStream(ctx context.Context) (Stream, error) {
	stream, err := qs.session.OpenStreamSync(ctx) //TODO: Check if we really want sync
	if err != nil {
		return nil, err
	}
	return &QUICStream{stream: stream}, nil
}

// QUICListener is a listener implementation for QUIC
type QUICListener struct {
	listener quic.Listener
}

// Close closes the QUICListener
func (l *QUICListener) Close() error {
	return l.listener.Close()
}

// Accept receives a connection for a given QUICListener
func (l *QUICListener) Accept(ctx context.Context) (Session, error) {
	session, err := l.listener.Accept(ctx)
	if err != nil {
		return nil, err
	}
	return &QUICSession{session: session}, nil
}

// QUICRPC RPC channel implementation for QUIC
type QUICRPC struct {
	address   string
	port      int
	tlsConfig *tls.Config
	config    *quic.Config
}

// Listen listens for incoming connections on QUICRPC
func (q *QUICRPC) Listen() (Listener, error) {
	// TODO: Check if IPV6 [] breaks this join logic
	listener, err := quic.ListenAddr(fmt.Sprintf("%s:%d", q.address, q.port), q.tlsConfig, q.config)
	if err != nil {
		return nil, err
	}
	return &QUICListener{listener: listener}, nil
}

// Connect connects on a server of QUICRPC
func (q *QUICRPC) Connect() (Session, error) {
	session, err := quic.DialAddr(fmt.Sprintf("%s:%d", q.address, q.port), q.tlsConfig, q.config)
	if err != nil {
		return nil, err
	}
	return &QUICSession{session: session}, nil
}

// NewQUICChannel Creates a channel using QUIC as transport layer
func NewQUICChannel(address string, port int, tlsConfig *tls.Config, quicConfig *quic.Config)RPC{
	return &QUICRPC{
		address:   address,
		port:      port,
		tlsConfig: tlsConfig,
		config:    quicConfig,
	}
}
