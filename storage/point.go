package storage

import (
	"log"
)

type Point struct {
	ID   string
	Body interface{}
}

func newPoint(PointId string, Body interface{}) *Point {
	return &Point{
		ID:   PointId,
		Body: Body,
	}
}

func _checkID(prefix, shieldID, pointID string) error {

	if shieldID == "" {
		return iotaToError(BadShieldID, prefix)
	}
	if pointID == "" {
		return iotaToError(BadPointID, prefix)
	}
	return nil
}

/*
	POINTS action
*/
// Set - adds or updates point into group
func Set(shieldID, pointID string, Body interface{}) error {
	return Singleton.Set(shieldID, pointID, Body)
}

func (s *Storage) Set(shieldID, pointID string, Body interface{}) error {

	if s.IsDebug {
		log.Printf("Set. shieldID: %s, pointID: %s\n", shieldID, pointID)
	}

	if err := _checkID("Set", shieldID, pointID); err != nil {
		return err
	}

	mesTo := newMessage(AddPoint, shieldID, pointID, Body)

	s.In <- mesTo
	mesFrom, ok := <-mesTo.Out
	if !ok {
		return iotaToError(InternalError, "AddPoint")
	}

	return iotaToError(mesFrom.Result, "AddPoint")
}

// Internal function. _setPoint - adds or updates point into group
func (s *Shield) _setPoint(pointID string, Body interface{}) uint {

	point := newPoint(pointID, Body)
	result := HookErrorPoint

	s.Lock()
	defer s.Unlock()

	shieldBody, pointBody, err := s._pointHookExe(AddPoint, point)

	if err == nil {
		result = Success
		s.Body = shieldBody
		point.Body = pointBody
		s.List[pointID] = point
	}

	return result
}

// Del - deletes one point if exists with channal interface
func Del(shieldID, pointID string) error {
	return Singleton.Del(shieldID, pointID)
}

// Del - deletes one point if exists with channal interface
func (s *Storage) Del(shieldID, pointID string) error {

	if s.IsDebug {
		log.Printf("Del. shieldID: %s, pointID: %s\n", shieldID, pointID)
	}

	if err := _checkID("Del", shieldID, pointID); err != nil {
		return err
	}

	mesTo := newMessage(DelPoint, shieldID, pointID)
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return iotaToError(InternalError, "Del")
	}
	return iotaToError(mesFrom.Result, "Del")
}

// Internal function. _delPoint - deletes one point
func (s *Shield) _delPoint(pointID string) uint {

	s.Lock()
	defer s.Unlock()

	if _, ok := s.List[pointID]; ok {
		delete(s.List, pointID)
	}

	return Success
}

// Get - returns one point if exists with channal interface
func Get(shieldID, pointID string) (interface{}, error) {
	return Singleton.Get(shieldID, pointID)
}

func GetMes(shieldID, pointID string) (interface{}, error) {
	return Singleton.GetMes(shieldID, pointID)
}

func (s *Storage) GetMes(shieldID, pointID string) (*Message, error) {

	if s.IsDebug {
		log.Printf("GetMes. shieldID: %s, pointID: %s\n", shieldID, pointID)
	}

	if err := _checkID("GetMes", shieldID, pointID); err != nil {
		return nil, err
	}

	mesTo := newMessage(GetPoint, shieldID, pointID)
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return nil, iotaToError(InternalError, "GetMes")
	}

	if mesFrom.Result != Success {
		return nil, iotaToError(mesFrom.Result, "GetMes")
	}

	return mesFrom, nil
}

func (s *Storage) Get(shieldID, pointID string) (interface{}, error) {

	if s.IsDebug {
		log.Printf("Get. shieldID: %s, pointID: %s\n", shieldID, pointID)
	}

	mes, err := s.GetMes(shieldID, pointID)
	if err != nil {
		return nil, err
	}

	return mes.Body, nil
}

// Internal function. _getPoint - returns one point if exists
func (s *Shield) _getPoint(pointID string) (point *Point, find bool) {

	s.RLock()
	point, find = s.List[pointID]
	s.RUnlock()

	return
}

// All - returns one points from the shield
func All(shieldID string) (map[string]AllResult, error) {
	return Singleton.All(shieldID)
}

func (s *Storage) All(shieldID string) (map[string]AllResult, error) {

	if s.IsDebug {
		log.Printf("All. shieldID %s\n", shieldID)
	}

	if shieldID == "" {
		return nil, iotaToError(BadShieldID, "All")
	}

	mesTo := newMessage(AllPoints, shieldID, "")
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return nil, iotaToError(InternalError, "All")
	}

	if mesFrom.Result != Success {
		return nil, iotaToError(mesFrom.Result, "All")
	}

	return mesFrom.All, nil
}

// Internal function. _getAllPoints - returns one points from the shield
func (s *Shield) _getAllPoints() (map[string]AllResult, uint) {

	s.RLock()
	defer s.RUnlock()

	out := map[string]AllResult{}
	for pointID, point := range s.List {
		_, pointBody, err := s._pointHookExe(AllPoints, point)
		out[pointID] = AllResult{
			Error: err,
			Body:  pointBody,
		}
	}

	return out, Success
}

func (s *Shield) _runEachFunc(f interface{}) uint {

	// func(interface{}, interface{}) (bool, interface{})
	fun, ok := f.(EachPointFunc)
	if !ok {
		return BadEachFunc
	}

	s.Lock()
	defer s.Unlock()

	for pointID, point := range s.List {
		isUpdate, body := fun(s.Body, point.Body)
		if isUpdate {
			s.List[pointID].Body = body
		}
	}

	return Success
}
