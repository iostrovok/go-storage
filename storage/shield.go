package storage

import (
	"log"
	"sync"
	"time"
)

type Shield struct {
	sync.RWMutex

	ID       string
	LastTime time.Time
	List     map[string]*Point
	In       chan *Message
	Body     interface{}
}

func newShields(shieldID string, Body interface{}) *Shield {
	s := &Shield{
		ID:       shieldID,
		LastTime: time.Now(),
		List:     map[string]*Point{},
		In:       make(chan *Message, 100),
		Body:     Body,
	}

	go s._Start()

	return s
}

func _checkOneID(prefix, shieldID string) error {

	if shieldID == "" {
		return iotaToError(BadShieldID, prefix)
	}

	return nil
}

/* Shields action */
// AddShield - adds new shield
func AddShield(shieldID string, Body interface{}) error {
	return Singleton.AddShield(shieldID, Body)
}

func (s *Storage) AddShield(shieldID string, Body interface{}) error {

	if s.IsDebug {
		log.Printf("AddShield. shieldID %s\n", shieldID)
	}

	if shieldID == "" {
		return iotaToError(BadShieldID, "AddShield")
	}

	mesTo := newMessage(AddGroup, shieldID, "", Body)
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return iotaToError(InternalError, "AddShield")
	}
	return iotaToError(mesFrom.Result, "AddShield")
}

func (s *Storage) _addShield(mes *Message) {

	s.Lock()

	if !s._shieldHookExe(AddGroup, mes) {
		return
	}

	_, find := s.Shields[mes.ShieldID]
	if find {
		mes.Result = ShieldExists
	} else {
		s.Shields[mes.ShieldID] = newShields(mes.ShieldID, mes.Body)
		mes.Result = Success
	}

	s.Unlock()
	mes.Out <- mes
}

// ExShield - checks existed shield
func ExShield(shieldID string) bool {
	return Singleton.ExShield(shieldID)
}

func (s *Storage) ExShield(shieldID string) bool {

	if s.IsDebug {
		log.Printf("ExShield. shieldID %s\n", shieldID)
	}

	if shieldID == "" {
		return false
	}

	s.RLock()
	_, find := s.Shields[shieldID]
	s.RUnlock()

	return find
}

// DelShield - deletes existed shield
func DelShield(shieldID string) error {
	return Singleton.DelShield(shieldID)
}

func (s *Storage) DelShield(shieldID string) error {

	if s.IsDebug {
		log.Printf("DelShield. shieldID %s\n", shieldID)
	}

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

	if !s._shieldHookExe(DelGroup, mes) {
		return
	}

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

// GetShield - returns data
func GetShield(shieldID string) (interface{}, error) {
	return Singleton.GetShield(shieldID)
}

func (s *Storage) GetShield(shieldID string) (interface{}, error) {

	if s.IsDebug {
		log.Printf("GetShield. shieldID: %s\n", shieldID)
	}

	if err := _checkOneID("GetMes", shieldID); err != nil {
		return nil, err
	}

	mesTo := newMessage(GetGroup, shieldID, "")
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return nil, iotaToError(InternalError, "GetShield")
	}

	if mesFrom.Result != Success {
		return nil, iotaToError(mesFrom.Result, "GetShield")
	}

	return mesFrom.Body, nil
}

// _getShield - returns exited shield
func (s *Storage) _getShield(mes *Message) (shield *Shield, find bool) {

	s.RLock()
	if !s._shieldHookExe(GetGroup, mes) {
		return
	}
	shield, find = s.Shields[mes.ShieldID]
	s.RUnlock()

	if !find {
		mes.Result = NotFoundShield
		mes.Out <- mes
	}

	return
}

func (s *Shield) _oneAct(mes *Message) {

	s.LastTime = time.Now()

	switch mes.Action {
	case AddPoint:
		mes.Result = s._setPoint(mes.PointId, mes.Body)
	case DelPoint:
		mes.Result = s._delPoint(mes.PointId)
	case GetPoint:
		mes.Body, mes.Result = s._getPoint(mes.PointId)
	case AllPoints:
		mes.All, mes.Result = s._getAllPoints()
	case UpdateTime:
		// Nothing
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
