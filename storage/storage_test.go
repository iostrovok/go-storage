package storage

import (
	//"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

func TestStorage(t *testing.T) {
	TestingT(t)
}

type StorageTestsSuite struct{}

var _ = Suite(&StorageTestsSuite{})

func (s *StorageTestsSuite) Test_New(c *C) {
	c.Skip("Not now")
	st := New()
	c.Assert(st, NotNil)
}

func (s *StorageTestsSuite) Test_newMessage(c *C) {
	c.Skip("Not now")
	m := newMessage(AddGroup, "ShieldID", "PointId")
	c.Assert(m, NotNil)
}

func (s *StorageTestsSuite) Test__addShield(c *C) {
	c.Skip("Not now")
	st := New()
	m := newMessage(AddGroup, "ShieldID", "")
	c.Assert(st.Shields["ShieldID"], IsNil)
	st._addShield(m)
	c.Assert(st.Shields["ShieldID"], NotNil)
}

func (s *StorageTestsSuite) Test_AddShield(c *C) {
	c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	c.Assert(st.Shields["ShieldID"], NotNil)
}

func (s *StorageTestsSuite) Test_AddShield_Empty(c *C) {
	c.Skip("Not now")
	st := New()
	err := st.AddShield("", "bla-bla-bla")
	c.Assert(err, NotNil)
}

func (s *StorageTestsSuite) Test__delShield(c *C) {
	c.Skip("Not now")
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
	c.Skip("Not now")
	st := New()
	err := st.DelShield("")
	c.Assert(err, NotNil)
}

func (s *StorageTestsSuite) Test_DelShield(c *C) {
	c.Skip("Not now")
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
	c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	c.Assert(st.Shields["ShieldID"], NotNil)

	st.Shields["ShieldID"]._setPoint("PointID", "sdsdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)
}

func (s *StorageTestsSuite) Test_Set(c *C) {
	c.Skip("Not now")
	st := New()

	c.Assert(st.Shields["ShieldID"], IsNil)

	err := st.AddShield("ShieldID", "bla-bla-bla")
	c.Assert(err, IsNil)

	st.Set("ShieldID", "PointID", "asdsad")

	c.Assert(st.Shields["ShieldID"].List["PointID"], NotNil)
}

func (s *StorageTestsSuite) Test__delPoint(c *C) {
	c.Skip("Not now")
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
	c.Skip("Not now")
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
	c.Skip("Not now")
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
	c.Skip("Not now")
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
	c.Skip("Not now")
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
