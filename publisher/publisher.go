package publisher

import (
	"bytes"
	"encoding/gob"
	"net"

	"github.com/Saf1u/pubsubshared/pubsubtypes"
)

func Publish(address string, message *pubsubtypes.Message) {
	con, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	data := make([]byte, 0, 1024)
	buffer := bytes.NewBuffer(data)
	encoder := gob.NewEncoder(buffer)
	err = encoder.Encode(*message)
	if err != nil {
		panic(err)
	}
	_, err = con.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}
