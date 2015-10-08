package storage

import (
	"log"
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

type HookShieldFunc func(interface{}) (interface{}, error)
type HookPointFunc func(interface{}, interface{}) (interface{}, interface{}, error)

type Storage struct {
	sync.RWMutex

	IsDebug bool

	PointHooks  map[uint][]HookPointFunc
	ShieldHooks map[uint][]HookShieldFunc

	Shields map[string]*Shield
	In      chan *Message
}

var Singleton *Storage = nil

func Debug(d ...bool) {
	Singleton.Debug(d...)
}

func (s *Storage) Debug(d ...bool) {
	if len(d) == 0 {
		s.IsDebug = true
	} else {
		s.IsDebug = d[0]
	}
}

func HookPoint(act uint, f HookPointFunc) {
	Singleton.HookPoint(act, f)
}

func HookShield(act uint, f HookShieldFunc) {
	Singleton.HookShield(act, f)
}

func (s *Storage) HookPoint(act uint, f HookPointFunc) {
	s.PointHooks[act] = append(s.PointHooks[act], f)
}

func (s *Storage) HookShield(act uint, f HookShieldFunc) {
	s.ShieldHooks[act] = append(s.ShieldHooks[act], f)
}

func New() *Storage {
	s := &Storage{
		Shields: map[string]*Shield{},
		In:      make(chan *Message, 1000),
	}

	s.PointHooks = map[uint][]HookPointFunc{}
	s.ShieldHooks = map[uint][]HookShieldFunc{}

	s.PointHooks[AddPoint] = []HookPointFunc{}
	s.PointHooks[DelPoint] = []HookPointFunc{}
	s.PointHooks[GetPoint] = []HookPointFunc{}

	s.ShieldHooks[AddGroup] = []HookShieldFunc{}
	s.ShieldHooks[DelGroup] = []HookShieldFunc{}
	s.ShieldHooks[GetGroup] = []HookShieldFunc{}

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

	if s.IsDebug {
		log.Printf("OneAct. mes.Action: %s\n", mes.Action)
	}

	switch mes.Action {
	case AddPoint, DelPoint, GetPoint:

		shield, find := s._getShield(mes)

		s._pointHookExe(mes.Action, shield, mes)

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
