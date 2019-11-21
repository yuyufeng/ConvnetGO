package convnetlib

import "encoding/json"

type Group struct {
	GroupName string
	GroupID   int
	Creator   int
	GroupDes  string
	NeedPass  bool
	Users     []User
}

func (group *Group) Json() string {
	data, err := json.Marshal(group)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func SliceClear(s *[]User) {
	*s = append([]User{})
}

func (group *Group) ClearUser() {
	SliceClear(&group.Users)
}

func (group *Group) Adduser(user *User) {
	group.Users = append(group.Users, *user)
}
