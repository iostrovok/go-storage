package storage

import (
	"sync"
)

const (
	AddPoint uint = iota
	AddGroup uint = iota
	DelPoint uint = iota
	DelGroup uint = iota
	GetPoint uint = iota
	GetGroup uint = iota
)

type Message struct {
	Action            uint
	Result            uint
	ShieldID, PointId string
	Out               chan *Message
	Body              interface{}
}

type Storage struct {
	sync.RWMutex

	Shields map[string]*Shield
	In      chan *Message
}

var Singleton *Storage = nil

func New() *Storage {
	s := &Storage{
		Shields: map[string]*Shield{},
		In:      make(chan *Message, 1000),
	}

	go s._start()
	return s
}

func StartSingleton() {
	if Singleton == nil {
		Singleton = New()
		go Singleton._start()
	}
}

// StopSingleton - use for test only
func StopSingleton() {
	if Singleton == nil {
		return
	}
	// TODO add deleting shields and points
	Singleton = nil
}

func newMessage(Action uint, ShieldID, PointId string, Body ...interface{}) *Message {
	mes := &Message{
		Action:   Action,
		ShieldID: ShieldID,
		PointId:  PointId,
		Out:      make(chan *Message, 1),
	}

	if len(Body) > 0 {
		mes.Body = Body
	}

	return mes
}

func (s *Storage) _start() {
	for {
		select {
		case mes, ok := <-s.In:
			if !ok {
				return
			}
			s.OneAct(mes)
		}
	}
}

func (s *Storage) OneAct(mes *Message) {

	switch mes.Action {
	case AddPoint, DelPoint, GetPoint:

		shield, find := s._getShield(mes)

		if find {
			shield.In <- mes
		}

	case AddGroup:
		s._addShield(mes)
	case DelGroup:
		s._delShield(mes)
	default:
		mes.Result = BadAction
		mes.Out <- mes
	}
}
