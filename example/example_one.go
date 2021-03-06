package main

import (
	"errors"
	"fmt"
	"log"
	"sync"

	st "storage"
)

type Group struct {
	sync.RWMutex

	ID    string
	Name  string
	Users map[string]*User
}

type User struct {
	ID   string
	Name string
}

func GetUser(group_id, user_id string) *User {
	m1, err := st.Get(group_id, user_id)
	if err != nil {
		log.Fatalln(err)
	}
	return m1.(*User)
}

// support function
func hook_func_shield(groupIn interface{}) (interface{}, error) {

	_, ok_g := groupIn.(*Group)
	if !ok_g {
		return groupIn, errors.New("No group type")
	}

	fmt.Printf("!!!!!!!!!!!!!!!!!hook_add_user START. groupIn: %+v\n", groupIn)

	return groupIn, nil
}

func hook_add_user(groupIn, userIn interface{}) (interface{}, interface{}, error) {

	fmt.Printf("hook_add_user START. groupIn: %+v\n", groupIn)
	fmt.Printf("hook_add_user START. userIn: %+v\n", userIn)

	group, ok_g := groupIn.(*Group)
	if !ok_g {
		return groupIn, userIn, errors.New("No group type")
	}
	user, ok_u := userIn.(*User)
	if !ok_u {
		return groupIn, userIn, errors.New("No user type")
	}

	group.Lock()
	group.Users[user.ID] = user
	group.Unlock()

	return group, user, nil
}

func main() {
	st.StartSingleton()
	st.Debug()
	st.HookShield(st.AddGroup, hook_func_shield)

	g := &Group{
		ID:    "my-group-ID",
		Name:  "my-group-NAME",
		Users: map[string]*User{},
	}

	st.AddShield("gr-1", g)

	user := &User{"1234567890", "Ivan"}
	st.Set("gr-1", "point-1", user)

	// Read user
	copyUser := GetUser("gr-1", "point-1")

	fmt.Printf("2. main. User: %+v\n", copyUser)

	// Read user
	copyGroup, err := st.GetShield("gr-1-wrong")
	fmt.Printf("3. main. err: %+v\n", err)
	fmt.Printf("4. main. copyGroup: %+v\n", copyGroup)
}
