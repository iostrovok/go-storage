package storage

import (
	"fmt"
	"log"
)

type HookShieldFunc func(interface{}) (interface{}, error)
type HookPointFunc func(interface{}, interface{}) (interface{}, interface{}, error)
type EachPointFunc func(interface{}, interface{}) (bool, interface{})

type ShieldHookMessage struct {
	Fun    HookPointFunc
	HookID string
	Act    uint
}

func HookShield(act uint, f HookShieldFunc) error {
	return Singleton.HookShield(act, f)
}

func (s *Storage) HookShield(act uint, f HookShieldFunc) error {
	switch act {
	case AddGroup, GetGroup:
		s.ShieldHooks[act] = append(s.ShieldHooks[act], f)
		return nil
	}

	return fmt.Errorf("Bad hook type in HookShield: %d", act)
}

func AddHookPoint(shieldID string, act uint, hookID string, f HookPointFunc) error {
	return Singleton.AddHookPoint(shieldID, act, hookID, f)
}

func (s *Storage) AddHookPoint(shieldID string, act uint, hookID string, f HookPointFunc) error {

	if s.IsDebug {
		log.Printf("AddHookPoint. shieldID: %s\n", shieldID)
	}

	shm := &ShieldHookMessage{
		Fun:    f,
		HookID: hookID,
		Act:    act,
	}

	mesTo := newMessage(AddHook, shieldID, "", shm)
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out

	if !ok {
		return iotaToError(InternalError, "AddHookPoint")
	}

	return iotaToError(mesFrom.Result, "AddHookPoint")
}

func DelHookPoint(shieldID string, act uint, hookID string) error {
	return Singleton.DelHookPoint(shieldID, act, hookID)
}

func (s *Storage) DelHookPoint(shieldID string, act uint, hookID string) error {

	if s.IsDebug {
		log.Printf("DelHookPoint. shieldID: %s\n", shieldID)
	}

	shm := &ShieldHookMessage{
		HookID: hookID,
		Act:    act,
	}

	mesTo := newMessage(DelHook, shieldID, "", shm)
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return iotaToError(InternalError, "AddHookPoint")
	}

	return iotaToError(mesFrom.Result, "AddHookPoint")
}

func (s *Storage) _shieldHookExe(Action uint, sh *Shield) (uint, interface{}) {

	if s.IsDebug {
		log.Printf("_shieldHookExe. Action: %d.\n", Action)
	}

	body := sh.Body
	var err error

	for _, f := range s.ShieldHooks[Action] {
		body, err = f(body)
		if err != nil {
			if s.IsDebug {
				log.Printf("_shieldHookExe. Action: %d. Error: %s.\n", err)
			}
			return HookErrorShield, nil
		}
	}

	return Success, body
}

func (s *Shield) _addHook(mes *Message) uint {

	shm, ok := mes.Body.(*ShieldHookMessage)
	if !ok {
		return HookErrorStruct
	}

	if shm.Fun == nil || shm.HookID == "" {
		return HookFuncBad
	}

	if shm.Act != AddPoint && shm.Act != GetPoint && shm.Act != AllPoints {
		return HookKeyBad
	}

	s.Lock()
	defer s.Unlock()
	s.Hooks[shm.Act][shm.HookID] = shm.Fun

	return Success
}

func (s *Shield) _delHook(mes *Message) uint {
	shm, ok := mes.Body.(*ShieldHookMessage)
	if !ok {
		return HookErrorStruct
	}

	if shm.Act != AddPoint && shm.Act != GetPoint && shm.Act != AllPoints {
		return HookKeyBad
	}

	s.Lock()
	defer s.Unlock()

	if _, ok := s.Hooks[shm.Act][shm.HookID]; ok {
		s.Hooks[shm.Act][shm.HookID] = nil
		delete(s.Hooks[shm.Act], shm.HookID)
	}

	return Success
}

func (s *Shield) _pointHookExe(Action uint, pt *Point) (shieldBody, pointBody interface{}, err error) {

	err = nil
	pointBody = pt.Body
	shieldBody = s.Body

	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("pkg: %v", r)
			}
		}
	}()

	for _, f := range s.Hooks[Action] {
		shieldBody, pointBody, err = f(shieldBody, pointBody)
		if err != nil {
			break
		}
	}

	return
}
