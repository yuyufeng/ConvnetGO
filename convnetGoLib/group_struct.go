package convnetlib

import (
	"encoding/json"
	"sync"
)

type Group struct {
	GroupName string
	GroupID   int
	Creator   int
	GroupDes  string
	NeedPass  bool
	Users     map[int64]*User
	sync.RWMutex
}

func (group *Group) Init() {
	if group.Users == nil {
		group.Users = make(map[int64]*User)
	}
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
	group.Lock()
	defer group.Unlock()
	group.Users[user.UserID] = user
}

func (group *Group) getUserByid(userid int64) (user *User) {
	group.Lock()
	defer group.Unlock()
	return group.Users[userid]
}

func (group *Group) removeUserByid(userid int64) {
	group.Lock()
	defer group.Unlock()
	delete(group.Users, userid)
}
