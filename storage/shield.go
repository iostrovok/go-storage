package storage

import (
	//"fmt"
	"sync"
	"time"
)

type Shield struct {
	sync.RWMutex

	ID       string
	LastTime time.Time
	List     map[string]*Point
	In       chan *Message
}

func newShields(shieldID string) *Shield {
	s := &Shield{
		ID:       shieldID,
		LastTime: time.Now(),
		List:     map[string]*Point{},
		In:       make(chan *Message, 100),
	}

	go s._Start()

	return s
}

/* Shields action */
// AddShield - adds new shield
func AddShield(shieldID string) error {
	return Singleton.AddShield(shieldID)
}

func (s *Storage) AddShield(shieldID string) error {

	if shieldID == "" {
		return iotaToError(BadShieldID, "AddShield")
	}

	mesTo := newMessage(AddGroup, shieldID, "")
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return iotaToError(InternalError, "AddShield")
	}
	return iotaToError(mesFrom.Result, "AddShield")
}

func (s *Storage) _addShield(mes *Message) {

	s.Lock()

	_, find := s.Shields[mes.ShieldID]
	if find {
		mes.Result = ShieldExists
	} else {
		s.Shields[mes.ShieldID] = newShields(mes.ShieldID)
		mes.Result = Success

	}

	s.Unlock()
	mes.Out <- mes
}

// DelShield - deletes existed shield
func DelShield(shieldID string) error {
	return Singleton.DelShield(shieldID)
}

func (s *Storage) DelShield(shieldID string) error {

	if shieldID == "" {
		return iotaToError(BadShieldID, "DelShield")
	}

	mesTo := newMessage(DelGroup, shieldID, "")
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return iotaToError(InternalError, "DelGroup")
	}
	return iotaToError(mesFrom.Result, "DelGroup")
}

func (s *Storage) _delShield(mes *Message) {

	s.Lock()
	_, find := s.Shields[mes.ShieldID]
	if !find {
		mes.Result = NotFoundShield
	} else {
		delete(s.Shields, mes.ShieldID)
		mes.Result = Success

	}
	s.Unlock()

	mes.Out <- mes
}

// _getShield - returns exited shield
func (s *Storage) _getShield(mes *Message) (shield *Shield, find bool) {

	s.RLock()
	shield, find = s.Shields[mes.ShieldID]
	s.RUnlock()

	if !find {
		mes.Result = NotFoundShield
		mes.Out <- mes
	}

	return
}

func (s *Shield) _oneAct(mes *Message) {

	switch mes.Action {
	case AddPoint:
		mes.Result = s._setPoint(mes.PointId, mes.Body)
	case DelPoint:
		mes.Result = s._delPoint(mes.PointId)
	case GetPoint:
		mes.Body, mes.Result = s._getPoint(mes.PointId)
	default:
		mes.Result = BadAction
	}
	mes.Out <- mes
}

func (s *Shield) _Start() {

	for {
		select {
		case mes, ok := <-s.In:

			if !ok {
				return
			}
			s._oneAct(mes)
		}
	}
}
