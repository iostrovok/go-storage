package storage

import ()

func (m Message) GetBody() interface{} {
	return m.Body
}

type AllResult struct {
	Error error
	Body  interface{}
}

type Message struct {
	Action            uint
	Result            uint
	ShieldID, PointId string
	Out               chan *Message
	Body              interface{}
	All               map[string]AllResult
}

func newMessage(Action uint, ShieldID, PointId string, Body ...interface{}) *Message {

	mes := &Message{
		Action:   Action,
		ShieldID: ShieldID,
		PointId:  PointId,
		Out:      make(chan *Message, 1),
		All:      map[string]AllResult{},
	}

	if len(Body) > 0 {
		mes.Body = Body[0]
	}

	return mes
}
