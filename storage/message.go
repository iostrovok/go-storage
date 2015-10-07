package storage

import (
	"fmt"
)

func (m Message) GetBody() interface{} {
	return m.Body
}

type Message struct {
	Action            uint
	Result            uint
	ShieldID, PointId string
	Out               chan *Message
	Body              interface{}
}

func newMessage(Action uint, ShieldID, PointId string, Body ...interface{}) *Message {

	mes := &Message{
		Action:   Action,
		ShieldID: ShieldID,
		PointId:  PointId,
		Out:      make(chan *Message, 1),
	}

	if len(Body) > 0 {
		mes.Body = Body[0]
	}

	fmt.Printf("newMessage: mes.Body: %+v\n", mes.Body)

	return mes
}
