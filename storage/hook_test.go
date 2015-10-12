package storage

import (
	"errors"
	. "gopkg.in/check.v1"
	"testing"
)

func StorageTestsHook(t *testing.T) {
	TestingT(t)
}

type StorageTestsSuite_Hook struct{}

var _ = Suite(&StorageTestsSuite_Hook{})

var hood_point_add string = "_added_hook_!"
var hood_shield_add string = "_added_hook_2"

// support function
func hook_func_shield(shieldBodyIn interface{}) (interface{}, error) {

	shieldBody, ok_shield := shieldBodyIn.(string)
	if !ok_shield {
		return shieldBodyIn, errors.New("No shield type")
	}

	shieldBody = shieldBody + hood_shield_add

	return shieldBody, nil
}

// support function
func hook_func_point(shieldBodyIn, pointBodyIn interface{}) (interface{}, interface{}, error) {

	shieldBody, ok_shield := shieldBodyIn.(string)
	if !ok_shield {
		return shieldBodyIn, pointBodyIn, errors.New("No shield type")
	}

	pointBody, ok_point := pointBodyIn.(string)
	if !ok_point {
		return shieldBodyIn, pointBodyIn, errors.New("No poimt type")
	}

	shieldBody = shieldBody + hood_point_add
	pointBody = pointBody + hood_point_add

	return shieldBody, pointBody, nil
}

// support function
func _check_string(line interface{}, val string) error {

	str, ok := line.(string)
	if !ok {
		return errors.New("No string type")
	}

	if str != val {
		return errors.New(str + " != " + val)
	}

	return nil
}

func (s *StorageTestsSuite_Hook) Test_HookPoint_HookShieldSimple(c *C) {
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	// Set hook
	err := st.HookShield(AddGroup, hook_func_shield)
	c.Assert(err, IsNil)

	c.Check(len(st.ShieldHooks[AddGroup]), Equals, 1)

}

func (s *StorageTestsSuite_Hook) Test_HookPoint_HookShield(c *C) {
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	// Set hook
	err := st.HookShield(AddGroup, hook_func_shield)
	c.Assert(err, IsNil)

	str := "bla-bla-bla"
	err = st.AddShield("ShieldID", str)
	c.Assert(err, IsNil)

	err = _check_string(st.Shields["ShieldID"].Body, str+hood_shield_add)
	c.Assert(err, IsNil)

}

func (s *StorageTestsSuite_Hook) Test_HookPoint_HookShieldError(c *C) {
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	// Set hook
	err := st.HookShield(AddPoint, hook_func_shield)
	c.Assert(err, NotNil)
	c.Check(len(st.ShieldHooks[AddGroup]), Equals, 0)
}

func (s *StorageTestsSuite_Hook) Test_HookPoint_HookShield_GetGroup(c *C) {
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	// Set hook
	errHook := st.HookShield(GetGroup, hook_func_shield)
	c.Assert(errHook, IsNil)

	str := "bla-bla-bla"
	err := st.AddShield("ShieldID", str)
	c.Assert(err, IsNil)

	err = _check_string(st.Shields["ShieldID"].Body, str)
	c.Assert(err, IsNil)

	body, err := st.GetShield("ShieldID")
	c.Assert(err, IsNil)

	err = _check_string(body, str+hood_shield_add)
	c.Assert(err, IsNil)

	err = _check_string(st.Shields["ShieldID"].Body, str)
	c.Assert(err, IsNil)
}

func (s *StorageTestsSuite_Hook) Test_HookPoint_HookShield_GetGroup_Error(c *C) {
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	var err error
	// Set hook
	errHook := st.HookShield(GetGroup, hook_func_shield)
	c.Assert(errHook, IsNil)

	str := "bla-bla-bla"
	err = st.AddShield("ShieldID", str)
	c.Assert(err, IsNil)

	err = _check_string(st.Shields["ShieldID"].Body, str)
	c.Assert(err, IsNil)

	body, err := st.GetShield("ShieldID-error")
	c.Assert(err, NotNil)
	c.Assert(body, IsNil)
}

func (s *StorageTestsSuite_Hook) Test_HookPoint_AddHookPoint(c *C) {
	var err error
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	str := "bla-bla-bla"
	err = st.AddShield("ShieldID", str+"_point_body")
	c.Assert(err, IsNil)

	// Set hook
	errHook := st.AddHookPoint("ShieldID", AddPoint, "super_hook", hook_func_point)
	c.Assert(errHook, IsNil)

	c.Assert(st.Shields["ShieldID"].Hooks[AddPoint]["super_hook"], NotNil)
}

func (s *StorageTestsSuite_Hook) Test_HookPoint_AddHookPoint_AddPoint(c *C) {
	var err error
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	str := "bla-bla-bla"
	err = st.AddShield("ShieldID", str+"_point_body")
	c.Assert(err, IsNil)

	// Set hook
	err = st.AddHookPoint("ShieldID", AddPoint, "super_hook", hook_func_point)
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID", str)
	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)

	errSh := _check_string(st.Shields["ShieldID"].Body, str+"_point_body"+hood_point_add)
	c.Assert(errSh, IsNil)

	errPt := _check_string(st.Shields["ShieldID"].List["PointID"].Body, str+hood_point_add)
	c.Assert(errPt, IsNil)
}

func (s *StorageTestsSuite_Hook) Test_HookPoint_AddHookPoint_ErrorActs(c *C) {
	var err error
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	str := "bla-bla-bla"
	err = st.AddShield("ShieldID", str+"_point_body")
	c.Assert(err, IsNil)

	Acts := []uint{GetPoint, AddPoint, AllPoints}
	for _, a := range Acts {
		// Set hook
		err = st.AddHookPoint("ShieldID", a, "super_hook", hook_func_point)
		c.Assert(err, IsNil)
	}

	noActs := []uint{AddGroup, DelPoint, DelGroup, GetGroup, Clean, UpdateTime, EachAct, DelHook, AddHook}
	for _, a := range noActs {
		// Set hook
		err = st.AddHookPoint("ShieldID", a, "super_hook", hook_func_point)
		c.Assert(err, NotNil)
	}

}

func (s *StorageTestsSuite_Hook) Test_HookPoint_AddHookPoint_GetPoint(c *C) {
	var err error
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	str := "bla-bla-bla"
	err = st.AddShield("ShieldID", str+"_point_body")
	c.Assert(err, IsNil)

	// Set hook
	err = st.AddHookPoint("ShieldID", GetPoint, "super_hook", hook_func_point)
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID", str)
	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)

	errSh := _check_string(st.Shields["ShieldID"].Body, str+"_point_body")
	c.Assert(errSh, IsNil)

	errPt := _check_string(st.Shields["ShieldID"].List["PointID"].Body, str)
	c.Assert(errPt, IsNil)

	body, errGet := st.Get("ShieldID", "PointID")
	c.Assert(errGet, IsNil)

	errBody := _check_string(body, str+hood_point_add)
	c.Assert(errBody, IsNil)
}

func (s *StorageTestsSuite_Hook) Test_HookPoint_AddHookPoint_AllPoints(c *C) {
	var err error
	//c.Skip("Not now")
	st := New()
	//st.Debug()

	str := "bla-bla-bla"
	err = st.AddShield("ShieldID", str+"_point_body")
	c.Assert(err, IsNil)

	// Set hook
	err = st.AddHookPoint("ShieldID", AllPoints, "super_hook", hook_func_point)
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID-1", str)
	c.Assert(st.Shields["ShieldID"].List["PointID-1"], NotNil)
	st.Set("ShieldID", "PointID-2", str)
	c.Assert(st.Shields["ShieldID"].List["PointID-2"], NotNil)

	errSh := _check_string(st.Shields["ShieldID"].Body, str+"_point_body")
	c.Assert(errSh, IsNil)

	errPt := _check_string(st.Shields["ShieldID"].List["PointID-1"].Body, str)
	c.Assert(errPt, IsNil)

	errPt = _check_string(st.Shields["ShieldID"].List["PointID-2"].Body, str)
	c.Assert(errPt, IsNil)

	body, errGet := st.All("ShieldID")
	c.Assert(errGet, IsNil)
	c.Assert(body, NotNil)

	for _, b := range body {

		c.Assert(b.Error, IsNil)
		errBody := _check_string(b.Body, str+hood_point_add)
		c.Assert(errBody, IsNil)
	}
}
