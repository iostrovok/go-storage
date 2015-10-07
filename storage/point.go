package storage

import (
//"fmt"
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

	s.Lock()
	defer s.Unlock()

	s.List[pointID] = newPoint(pointID, Body)

	return Success
}

// Del - deletes one point if exists with channal interface
func Del(shieldID, pointID string) error {
	return Singleton.Del(shieldID, pointID)
}

// Del - deletes one point if exists with channal interface
func (s *Storage) Del(shieldID, pointID string) error {

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

func (s *Storage) Get(shieldID, pointID string) (interface{}, error) {

	if err := _checkID("Get", shieldID, pointID); err != nil {
		return nil, err
	}

	mesTo := newMessage(GetPoint, shieldID, pointID)
	s.In <- mesTo

	mesFrom, ok := <-mesTo.Out
	if !ok {
		return nil, iotaToError(InternalError, "Get")
	}

	if mesFrom.Result != Success {
		return nil, iotaToError(mesFrom.Result, "Get")
	}

	return mesFrom.Result, nil
}

// Internal function. _getPoint - returns one point if exists
func (s *Shield) _getPoint(pointID string) (interface{}, uint) {

	s.Lock()
	defer s.Unlock()

	point, ok := s.List[pointID]
	if !ok {
		return nil, NotFoundPoint
	}

	return point.Body, Success
}
