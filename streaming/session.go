package streaming

import (
	"bufio"
	"crypto/tls"

	"github.com/belmegatron/gofair/streaming/models"
)

const (
	// Error Codes for Authenticating
	failure = "FAILURE"
)

type Session struct {
	conn         *TLSConnection
	scanner      *bufio.Scanner
	eventHandler *EventHandler
	stop         chan int
	channels     *StreamChannels
}

func NewSession(destination string, certs *tls.Certificate, appKey string, sessionToken string, channels *StreamChannels) (*Session, error) {
	session := new(Session)
	TLSConnection, err := NewTLSConnection(destination, certs)
	if err != nil {
		return nil, err
	}

	session.conn = TLSConnection

	// Wrap the underlying connection with our byte scanner, this will read in bytes until a CRLF is encountered
	session.scanner = bufio.NewScanner(TLSConnection.conn)
	session.scanner.Split(scanCRLF)

	// Pass a pointer to our StreamChannels struct which are used for piping data back to the main goroutine
	session.channels = channels

	err = session.authenticate(appKey, sessionToken)
	if err != nil {
		return nil, err
	}

	go session.readPump()
	go session.writePump()

	return session, nil
}

func (session *Session) Stop() {
	// Stop the readPump/writePump goroutines
	session.stop <- 1
	// Terminate TLS connection to stream endpoint
	session.conn.Stop()
}

func (session *Session) authenticate(appKey string, sessionToken string) error {

	if session.conn == nil {
		return &NoConnectionError{}
	}

	authenticationMessage := &models.AuthenticationMessage{AppKey: appKey, Session: sessionToken}
	authenticationMessage.SetID(session.conn.ID)

	b, err := authenticationMessage.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = session.write(b)
	if err != nil {
		return err
	}

	buf, err := session.read()
	if err != nil {
		return err
	}

	statusMessage := new(models.StatusMessage)
	err = statusMessage.UnmarshalJSON(buf)
	if err != nil {
		return err
	}

	if statusMessage.StatusCode == failure {
		return &AuthenticationError{}
	}

	return nil
}

func (session *Session) write(b []byte) (int, error) {
	// Every message is in json & terminated with a line feed (CRLF)
	b = addCRLF(b)
	return session.conn.Write(b)
}

func (session *Session) read() ([]byte, error) {

	session.scanner.Scan()

	if err := session.scanner.Err(); err != nil {
		return []byte{}, err
	}

	return session.scanner.Bytes(), nil
}

func (session *Session) readPump() {

	if session.conn == nil {
		err := new(NoConnectionError)
		session.channels.Err <- err
		return
	}

	for {
		select {

		case <-session.stop:
			return

		default:
			buf, err := session.read()

			if err != nil {
				session.channels.Err <- err
				return
			}

			op, err := getOp(buf)
			if err != nil {
				session.channels.Err <- err
				return
			}

			session.eventHandler.onData(op, buf)
		}
	}
}

func (session *Session) writePump() {
	for {
		select {

		case <-session.stop:
			return

		case marketSubscriptionMessage := <-session.channels.marketSubscriptionRequest:
			b, err := marketSubscriptionMessage.MarshalJSON()
			if err != nil {
				session.channels.Err <- err
				return
			}

			session.write(b)

		case orderSubscriptionMessage := <-session.channels.orderSubscriptionRequest:

			b, err := orderSubscriptionMessage.MarshalJSON()
			if err != nil {
				session.channels.Err <- err
				return
			}

			session.write(b)
		}
	}
}