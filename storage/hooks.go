package storage

import (
	"log"
)

type HookShieldFunc func(interface{}) (interface{}, error)
type HookPointFunc func(interface{}, interface{}) (interface{}, interface{}, error)

func HookPoint(act uint, f HookPointFunc) {
	Singleton.HookPoint(act, f)
}

func (s *Storage) HookPoint(act uint, f HookPointFunc) {
	s.PointHooks[act] = append(s.PointHooks[act], f)
}

func HookShield(act uint, f HookShieldFunc) {
	Singleton.HookShield(act, f)
}

func (s *Storage) HookShield(act uint, f HookShieldFunc) {
	s.ShieldHooks[act] = append(s.ShieldHooks[act], f)
}

func (s *Storage) _shieldHookExe(Action uint, mes *Message) bool {

	if s.IsDebug {
		log.Printf("_shieldHookExe. Action: %d.\n", Action)
	}

	for _, f := range s.ShieldHooks[Action] {
		res, err := f(mes.Body)
		if err != nil {
			mes.Result = HookErrorShield
			mes.Out <- mes
			return false
		}
		mes.Body = res
	}
	return true
}

func (s *Storage) _pointHookExe(Action uint, sh *Shield, mes *Message) bool {

	if s.IsDebug {
		log.Printf("_pointHookExe. Action: %d.\n", Action)
	}

	for _, f := range s.PointHooks[Action] {
		shBody, mesBody, err := f(sh.Body, mes.Body)
		if err != nil {
			mes.Result = HookErrorPoint
			mes.Out <- mes
			return false
		}
		sh.Body = shBody
		mes.Body = mesBody
	}
	return true
}
