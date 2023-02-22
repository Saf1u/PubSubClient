package subscriber

import (
	"bytes"
	"encoding/gob"
	"net"

	"github.com/Saf1u/pubsubshared/pubsubtypes"
)

type SuscriberClient struct {
	con    net.Conn
	id     int
	buffer []byte
	topic  string
}

func RegisterSuscriber(address string, topic string) *SuscriberClient {
	regMessage := pubsubtypes.Message{Type: pubsubtypes.REGISTER_CONN, Topic: topic}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	data := make([]byte, 0, 100)
	bufferr := bytes.NewBuffer(data)
	encoder := gob.NewEncoder(bufferr)
	err = encoder.Encode(regMessage)
	toSend := bufferr.Bytes()
	length := make([]byte, 1)
	length[0] = byte(bufferr.Len())
	length = append(length, toSend...)
	if err != nil {
		panic(err)
	}

	subs := &SuscriberClient{con: conn, buffer: make([]byte, 1024), topic: topic}
	_, err = conn.Write(length)
	if err != nil {
		panic(err)
	}
	regMessage = *subs.read()
	subs.id = regMessage.Id
	return subs

}
func (s *SuscriberClient) read() *pubsubtypes.Message {
	s.con.Read(s.buffer[0:1])
	s.con.Read(s.buffer[1 : s.buffer[0]+1])
	reader := bytes.NewReader(s.buffer[1 : s.buffer[0]+1])
	decoder := gob.NewDecoder(reader)
	msg := &pubsubtypes.Message{}
	err := decoder.Decode(msg)
	if err != nil {
		panic(err)
	}
	return msg
}

func (s *SuscriberClient) Close() {
	msg := &pubsubtypes.Message{
		Id:    s.id,
		Type:  pubsubtypes.CLOSE_CONN,
		Topic: s.topic,
	}
	data := make([]byte, 0, 1024)
	buffer := bytes.NewBuffer(data)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(*msg)
	if err != nil {
		panic(err)
	}
	_, err = s.con.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (s *SuscriberClient) Read() *pubsubtypes.Message {
	msg := s.read()
	return msg
}
