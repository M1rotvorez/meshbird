package protocol

import (
	"bytes"
	"fmt"
	"github.com/gophergala2016/meshbird/secure"
	"io"
)

var (
	magicKey = []byte{'M', 'E', 'S', 'H', 'B', 'I', 'R', 'D'}
)

type (
	HandshakeMessage []byte
)

func IsMagicValid(data []byte) bool {
	logger.Printf("Trying to check magic (%v) : %v", magicKey, data)
	return bytes.HasPrefix(data, magicKey)
}

func NewHandshakePacket(sessionKey []byte, networkSecret *secure.NetworkSecret) *Packet {
	sessionKey = append(magicKey, sessionKey...)
	data := networkSecret.Encode(sessionKey)

	body := Body{
		Type: TypeHandshake,
		Msg:  HandshakeMessage(data),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m HandshakeMessage) Len() uint16 {
	return uint16(len(m))
}

func (m HandshakeMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}

func (m HandshakeMessage) Bytes() []byte {
	return []byte(m)
}

func (m HandshakeMessage) SessionKey() []byte {
	return m[len(magicKey):]
}

func ReadDecodeHandshake(r io.Reader) (HandshakeMessage, error) {
	logger.Printf("Trying to read Handshake message...")

	okPack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.Printf("Unable to decode package: %s", errDecode)
		return nil, fmt.Errorf("Error on read Handshake package: %v", errDecode)
	}

	if okPack.Data.Type != TypeHandshake {
		return nil, fmt.Errorf("Got non Handshake message: %+v", okPack)
	}

	logger.Printf("Readed Handshake: %+v", okPack.Data.Msg)

	return okPack.Data.Msg.(HandshakeMessage), nil
}

func WriteEncodeHandshake(w io.Writer, sessionKey []byte, networkSecret *secure.NetworkSecret) (err error) {
	logger.Printf("Trying to write Handshake message...")
	if err = EncodeAndWrite(w, NewHandshakePacket(sessionKey, networkSecret)); err != nil {
		err = fmt.Errorf("Error on write Handshake message: %v", err)
	}
	return
}
