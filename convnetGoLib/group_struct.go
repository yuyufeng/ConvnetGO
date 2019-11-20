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

func (group *Group) UserJson() string {
	data, err := json.Marshal(group.Users)
	if err != nil {
		panic(err)
	}
	return string(data)
}
