package storage

import (
	"log"
	"sync"
	"time"
)

const (
	AddPoint uint = iota
	AddGroup uint = iota
	DelPoint uint = iota
	DelGroup uint = iota
	GetPoint uint = iota
	GetGroup uint = iota

	AllPoints uint = iota
	Clean     uint = iota

	UpdateTime uint = iota
)

type Storage struct {
	sync.RWMutex

	IsDebug bool

	PointHooks  map[uint][]HookPointFunc
	ShieldHooks map[uint][]HookShieldFunc

	Shields map[string]*Shield
	In      chan *Message
	TimerCh chan time.Duration

	ShieldTTL time.Duration
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

func New() *Storage {
	s := &Storage{
		Shields:   map[string]*Shield{},
		In:        make(chan *Message, 1000),
		TimerCh:   make(chan time.Duration, 10),
		ShieldTTL: time.Minute,
	}

	s.PointHooks = map[uint][]HookPointFunc{}
	s.ShieldHooks = map[uint][]HookShieldFunc{}

	s.PointHooks[AllPoints] = []HookPointFunc{}
	s.PointHooks[AddPoint] = []HookPointFunc{}
	s.PointHooks[DelPoint] = []HookPointFunc{}
	s.PointHooks[GetPoint] = []HookPointFunc{}

	s.ShieldHooks[AllPoints] = []HookShieldFunc{}
	s.ShieldHooks[AddGroup] = []HookShieldFunc{}
	s.ShieldHooks[DelGroup] = []HookShieldFunc{}
	s.ShieldHooks[GetGroup] = []HookShieldFunc{}

	go s._start()
	go s._check_ttl()
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

func (s *Storage) _cleanShield() {

	if s.IsDebug {
		log.Printf("_cleanShield. %s\n")
	}

	n := time.Now().Add(-1 * s.ShieldTTL)
	newShields := map[string]*Shield{}

	s.Lock()

	for key, shield := range s.Shields {
		if n.After(shield.LastTime) {
			continue
		}
		newShields[key] = shield
	}

	if s.IsDebug {
		log.Printf("_cleanShield result: was: %d, now: %d\n", len(s.Shields), len(newShields))
	}

	s.Shields = newShields

	s.Unlock()
}

func (s *Storage) OneAct(mes *Message) {

	if s.IsDebug {
		log.Printf("OneAct. mes.Action: %s\n", mes.Action)
	}

	switch mes.Action {
	case AddPoint, DelPoint, GetPoint, AllPoints:

		shield, find := s._getShield(mes)

		if find {
			s._pointHookExe(mes.Action, shield, mes)
			shield.In <- mes
		}

	case GetGroup:

		shield, find := s._getShield(mes)

		if find {
			mes.Body = shield.Body
			mes.Result = Success
			mes.Out <- mes

			shield.In <- newMessage(UpdateTime, "", "", "")
		}

	case AddGroup:
		s._addShield(mes)
	case DelGroup:
		s._delShield(mes)
	case Clean:
		s._cleanShield()
	default:
		mes.Result = BadAction
		mes.Out <- mes
	}
}
