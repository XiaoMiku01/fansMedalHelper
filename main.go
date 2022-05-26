package main

import (
	"MedalHelper/service"
	"MedalHelper/util"
	"strconv"
	"strings"
)

func InitUsers() []service.User {
	users := make([]service.User, 0, 1)
	for _, userInfo := range util.GlobalConfig.UserList {
		banId := make([]int, 0)
		if userInfo.BannedUid != "" {
			banIdStr := strings.Split(userInfo.BannedUid, ",")
			for _, str := range banIdStr {
				id, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					continue
				}
				banId = append(banId, int(id))
			}
		}
		users = append(users, service.NewUser(userInfo.AccessKey, banId))
	}
	return users
}

func main() {
	users := InitUsers()
	for _, user := range users {
		user.Init()
		user.Start()
	}
}
