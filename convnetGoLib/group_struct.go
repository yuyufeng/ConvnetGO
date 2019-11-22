package convnetlib

import (
	"encoding/json"
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
	GroupID   int64
	Creator   int
	GroupDes  string
	NeedPass  bool
	Users     map[int64]*User
	sync.RWMutex
}

func NewGroup() *Group {
	var result = new(Group)
	if result.Users == nil {
		result.Users = make(map[int64]*User)
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
	group.Users = make(map[int64]*User)
}

func (group *Group) Adduser(user *User) {

	if user.UserID == client.MyUserid {
		return
	}

	group.Lock()
	defer group.Unlock()
	group.Users[user.UserID] = user
}

func (group *Group) GetUserByid(userid int64) (user *User) {
	group.Lock()
	defer group.Unlock()
	return group.Users[userid]
}

func (group *Group) RemoveUserByid(userid int64) {
	group.Lock()
	defer group.Unlock()
	delete(group.Users, userid)
}
