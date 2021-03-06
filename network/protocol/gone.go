package protocol

import (
	"fmt"
	"io"
)

type (
	GoneMessage []byte
)

func NewGoneMessage() *Packet {
	body := Body{
		Type: TypeGone,
		Msg:  GoneMessage([]byte{}),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m GoneMessage) Len() uint16 {
	return uint16(len(m))
}

func (m GoneMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}

func ReadDecodeGone(r io.Reader) (GoneMessage, error) {
	logger.Printf("Trying to read Gone message...")

	gonePack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.Printf("Unable to decode package: %s", errDecode)
		return nil, fmt.Errorf("Error on read Gone package: %v", errDecode)
	}

	if gonePack.Data.Type != TypeGone {
		return nil, fmt.Errorf("Got non Gone message: %+v", gonePack)
	}

	logger.Printf("Readed Gone: %+v", gonePack.Data.Msg)

	return gonePack.Data.Msg.(GoneMessage), nil
}

func WriteEncodeGone(w io.Writer) (err error) {
	logger.Printf("Trying to write Gone message...")
	if err = EncodeAndWrite(w, NewGoneMessage()); err != nil {
		err = fmt.Errorf("Error on write Gone message: %v", err)
	}
	return
}
