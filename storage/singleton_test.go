package storage

import (
	//"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

func TestStorageSingleton(t *testing.T) {
	TestingT(t)
}

type StorageTestsSuite_Singleton struct{}

var _ = Suite(&StorageTestsSuite_Singleton{})

func (s *StorageTestsSuite_Singleton) Test_StartSingleton(c *C) {
	// c.Skip("Not now")
	StartSingleton()
	c.Assert(Singleton, NotNil)
}

func (s *StorageTestsSuite_Singleton) Test_StopSingleton(c *C) {
	//c.Skip("Not now")
	StartSingleton()
	c.Assert(Singleton, NotNil)
	StopSingleton()
	c.Assert(Singleton, IsNil)
}

func (s *StorageTestsSuite_Singleton) Test_AddShield_Singleton(c *C) {
	//c.Skip("Not now")

	StartSingleton()
	c.Assert(Singleton.Shields["ShieldID"], IsNil)
	err := AddShield("ShieldID")
	c.Assert(err, IsNil)
	c.Assert(Singleton.Shields["ShieldID"], NotNil)
}

func (s *StorageTestsSuite_Singleton) Test_AddShield_Empty_Singleton(c *C) {
	//c.Skip("Not now")
	StartSingleton()
	err := AddShield("")
	c.Assert(err, NotNil)
}

func (s *StorageTestsSuite_Singleton) Test_DelShield_Empty_Singleton(c *C) {
	//c.Skip("Not now")
	StartSingleton()
	err := DelShield("")
	c.Assert(err, NotNil)
}

func (s *StorageTestsSuite_Singleton) Test_DelShield_Singleton(c *C) {
	//c.Skip("Not now")

	StopSingleton()
	StartSingleton()

	c.Assert(Singleton.Shields["ShieldID"], IsNil)

	err := AddShield("ShieldID")
	c.Assert(err, IsNil)

	c.Assert(Singleton.Shields["ShieldID"], NotNil)
	err = DelShield("ShieldID")
	c.Assert(err, IsNil)

	c.Assert(Singleton.Shields["ShieldID"], IsNil)
}

func (s *StorageTestsSuite_Singleton) Test_Set_Singleton(c *C) {
	//c.Skip("Not now")
	StopSingleton()
	StartSingleton()

	c.Assert(Singleton.Shields["ShieldID"], IsNil)

	err := AddShield("ShieldID")
	c.Assert(err, IsNil)

	Set("ShieldID", "PointID", "asdsad")

	c.Assert(Singleton.Shields["ShieldID"].List["PointID"], NotNil)
}

func (s *StorageTestsSuite_Singleton) Test_Del_Singleton(c *C) {
	//c.Skip("Not now")
	StopSingleton()
	StartSingleton()

	c.Assert(Singleton.Shields["ShieldID"], IsNil)

	err := AddShield("ShieldID")
	c.Assert(err, IsNil)

	Set("ShieldID", "PointID", "asdsad")

	c.Assert(Singleton.Shields["ShieldID"].List["PointID"], NotNil)

	res := Del("ShieldID", "PointID")
	c.Assert(res, IsNil)
	c.Assert(Singleton.Shields["ShieldID"].List["PointID"], IsNil)
}

func (s *StorageTestsSuite_Singleton) Test_Get_Singleton(c *C) {
	//c.Skip("Not now")
	StopSingleton()
	StartSingleton()

	c.Assert(Singleton.Shields["ShieldID"], IsNil)

	err := AddShield("ShieldID")
	c.Assert(err, IsNil)

	Set("ShieldID", "PointID", "asdsad")

	c.Assert(Singleton.Shields["ShieldID"].List["PointID"], NotNil)

	body, res := Get("ShieldID", "PointID")
	c.Assert(Singleton.Shields["ShieldID"].List["PointID"], NotNil)
	c.Assert(body, NotNil)
	c.Assert(res, IsNil)
}

func (s *StorageTestsSuite_Singleton) Test_Get_Error_Singleton(c *C) {
	//c.Skip("Not now")
	StopSingleton()
	StartSingleton()

	c.Assert(Singleton.Shields["ShieldID"], IsNil)

	err := AddShield("ShieldID")
	c.Assert(err, IsNil)

	Set("ShieldID", "PointID", "asdsad")

	c.Assert(Singleton.Shields["ShieldID"].List["PointID"], NotNil)

	body, res := Get("ShieldID", "PointID-error")
	c.Assert(Singleton.Shields["ShieldID"].List["PointID"], NotNil)
	c.Assert(res, NotNil)
	c.Assert(body, IsNil)
}
