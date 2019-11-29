package convnetlib

import (
	"bytes"
	"encoding/json"
	"net"
	"sync"
)

type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type WaitGroup struct {
	noCopy noCopy
	state1 [3]uint32
}

type Group struct {
	noCopy    noCopy
	GroupName string
	GroupID   int
	Creator   int
	GroupDes  string
	NeedPass  bool
	Users     map[int]*User
	sync.RWMutex
}

func NewGroup() *Group {
	var result = new(Group)
	if result.Users == nil {
		result.Users = make(map[int]*User)
	}
	return result
}

func (group *Group) Json() string {
	group.Lock()
	defer group.Unlock()
	data, err := json.Marshal(group)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (group *Group) ClearUser() {
	group.Lock()
	defer group.Unlock()
	group.Users = make(map[int]*User)
}

func (group *Group) Adduser(user *User) {

	if user.UserID == client.MyUserid {
		return
	}

	group.Lock()
	defer group.Unlock()
	group.Users[user.UserID] = user
}

func (group *Group) GetUserByid(userid int) (user *User) {
	group.Lock()
	defer group.Unlock()
	return group.Users[userid]
}

func (group *Group) RemoveUserByid(userid int) {
	group.Lock()
	defer group.Unlock()
	delete(group.Users, userid)
}

func GetUserByMac(mac net.HardwareAddr) (user *User) {
	group := client.g_AllUser
	group.Lock()
	defer group.Unlock()
	for i, v := range group.Users {
		if bytes.Equal(v.MacAddress, mac) {
			return group.Users[i]
		}
	}
	return nil
}
