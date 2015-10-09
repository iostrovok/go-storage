package storage

import (
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	TestingT(t)
}

type StorageTestsSuite struct{}

var _ = Suite(&StorageTestsSuite{})

func (s *StorageTestsSuite) Test_New(c *C) {
	//c.Skip("Not now")
	st := New()
	c.Assert(st, NotNil)
}

func (s *StorageTestsSuite) Test_newMessage(c *C) {
	//c.Skip("Not now")
	m := newMessage(AddGroup, "ShieldID", "PointId")
	c.Assert(m, NotNil)
}

func (s *StorageTestsSuite) Test__addShield(c *C) {
	//c.Skip("Not now")
	st := New()
	m := newMessage(AddGroup, "ShieldID", "")
	c.Assert(st.Shields["ShieldID"], IsNil)
	st._addShield(m)
	c.Assert(st.Shields["ShieldID"], NotNil)
}

func (s *StorageTestsSuite) Test_AddShield(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	c.Assert(st.Shields["ShieldID"], NotNil)
}

func (s *StorageTestsSuite) Test_AddShield_Empty(c *C) {
	//c.Skip("Not now")
	st := New()
	err := st.AddShield("", "bla-bla-bla")
	c.Assert(err, NotNil)
}

func (s *StorageTestsSuite) Test__delShield(c *C) {
	//c.Skip("Not now")
	st := New()
	m := newMessage(AddGroup, "ShieldID", "")
	c.Assert(st.Shields["ShieldID"], IsNil)
	st._addShield(m)
	c.Assert(st.Shields["ShieldID"], NotNil)

	m2 := newMessage(DelGroup, "ShieldID", "")
	st._delShield(m2)
	c.Assert(st.Shields["ShieldID"], IsNil)
}

func (s *StorageTestsSuite) Test_DelShield_Empty(c *C) {
	//c.Skip("Not now")
	st := New()
	err := st.DelShield("")
	c.Assert(err, NotNil)
}

func (s *StorageTestsSuite) Test_DelShield(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	c.Assert(st.Shields["ShieldID"], NotNil)
	err = st.DelShield("ShieldID")
	c.Assert(err, IsNil)

	c.Assert(st.Shields["ShieldID"], IsNil)
}

func (s *StorageTestsSuite) Test__setPoint(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	c.Assert(st.Shields["ShieldID"], NotNil)

	st.Shields["ShieldID"]._setPoint("PointID", "sdsdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)
}

func (s *StorageTestsSuite) Test_Set(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID", "asdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)
}

func (s *StorageTestsSuite) Test__delPoint(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	c.Assert(st.Shields["ShieldID"], NotNil)

	st.Shields["ShieldID"]._setPoint("PointID", "sdsdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)

	res := st.Shields["ShieldID"]._delPoint("PointID")
	c.Assert(res, Equals, Success)
	c.Assert(st.Shields["ShieldID"].List["PointID"], IsNil)
}

func (s *StorageTestsSuite) Test_Del(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID", "asdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)

	res := st.Del("ShieldID", "PointID")
	c.Assert(res, IsNil)
	c.Assert(st.Shields["ShieldID"].List["PointID"], IsNil)
}

func (s *StorageTestsSuite) Test__getPoint(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	c.Assert(st.Shields["ShieldID"], NotNil)

	st.Shields["ShieldID"]._setPoint("PointID", "sdsdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)

	body, res := st.Shields["ShieldID"]._getPoint("PointID")
	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)
	c.Assert(body, NotNil)
	c.Assert(res, Equals, Success)
}

func (s *StorageTestsSuite) Test_Get(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID", "asdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)

	body, res := st.Get("ShieldID", "PointID")
	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)
	c.Assert(body, NotNil)
	c.Assert(res, IsNil)
}

func (s *StorageTestsSuite) Test_Get_Error(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID", "asdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)

	body, res := st.Get("ShieldID", "PointID-error")
	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)
	c.Assert(res, NotNil)
	c.Assert(body, IsNil)
}

func (s *StorageTestsSuite) Test_All(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla-1")
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID-1", "bla-bla-bla-1")
	st.Set("ShieldID", "PointID-2", "bla-bla-bla-2")
	st.Set("ShieldID", "PointID-3", "bla-bla-bla-3")

	c.Assert(len(st.Shields["ShieldID"].List), Equals, 3)

	all, res := st.All("ShieldID")

	c.Assert(all, NotNil)
	c.Assert(res, IsNil)
	c.Assert(len(all), Equals, 3)

	str, ok := all["PointID-2"].(string)
	c.Assert(ok, Equals, true)
	c.Assert(str, Equals, "bla-bla-bla-2")
}

func (s *StorageTestsSuite) Test_GetShield(c *C) {
	//c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID", "asdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)

	body, res := st.GetShield("ShieldID")
	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)
	c.Assert(body, NotNil)
	c.Assert(res, IsNil)

	str, ok := body.(string)
	c.Assert(ok, Equals, true)
	c.Assert(str, Equals, "bla-bla-bla")
}

func (s *StorageTestsSuite) Test_SetShieldTTL(c *C) {
	//c.Skip("Not now")
	st := New()
	st.SetShieldTTL(time.Minute)

	c.Assert(st.Shields["ShieldID"], IsNil)

	st.AddShield("ShieldID", "bla-bla-bla-1")
	st.AddShield("ShieldID-2", "bla-bla-bla-1")

	st.Set("ShieldID", "PointID-1", "bla-bla-bla-1")
	st.Set("ShieldID", "PointID-2", "bla-bla-bla-2")
	st.Set("ShieldID-2", "PointID-3", "bla-bla-bla-3")

	c.Assert(len(st.Shields), Equals, 2)

	st.SetShieldTTL(time.Millisecond * 500)
	time.Sleep(time.Millisecond * 1000)

	fmt.Printf("%+v\n", st.Shields)

	c.Assert(len(st.Shields), Equals, 0)
}
