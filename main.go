package main

import (
	"MedalHelper/service"
	"MedalHelper/util"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func usage() {
	fmt.Print(`Usage: main.go [command]

command:
    login   login bili account and get access key
`)
}

func initUsers() []service.User {
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

func exec() {
	users := initUsers()
	wg := sync.WaitGroup{}
	for _, user := range users {
		if status := user.Init(); status {
			wg.Add(1)
			go user.Start(wg)
		}
	}
	wg.Wait()
}

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "login" {
			util.LoginBili()
		} else {
			usage()
		}
		return
	}

	exec()
}
