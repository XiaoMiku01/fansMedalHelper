package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	gjson "github.com/tidwall/gjson"
)

var accessKey string

func getQRcode() (string, string) {
	api := "http://passport.bilibili.com/x/passport-tv-login/qrcode/auth_code"
	data := map[string]string{
		"local_id": "0",
		"ts":       GetTimestamp(),
	}
	Signature(&data)
	data_string := strings.NewReader(Map2string(data))
	client := http.Client{}
	req, _ := http.NewRequest("POST", api, data_string)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	code := gjson.Parse(string(body)).Get("code").Int()
	if code == 0 {
		qrcode_url := gjson.Parse(string(body)).Get("data.url").String()
		auth_code := gjson.Parse(string(body)).Get("data.auth_code").String()
		return qrcode_url, auth_code
	} else {
		// FIXME: handle error here
		panic("getQRcode error")
	}
}

func verifyLogin(auth_code string) {
	for {
		time.Sleep(time.Second * 3)
		api := "http://passport.bilibili.com/x/passport-tv-login/qrcode/poll"
		data := map[string]string{
			"auth_code": auth_code,
			"local_id":  "0",
			"ts":        GetTimestamp(),
		}
		Signature(&data)
		data_string := strings.NewReader(Map2string(data))
		client := http.Client{}
		req, _ := http.NewRequest("POST", api, data_string)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		code := gjson.Parse(string(body)).Get("code").Int()
		accessKey = gjson.Parse(string(body)).Get("data.access_token").String()
		if code == 0 {
			fmt.Println("登录成功")
			fmt.Println("access_key:", string(accessKey))
			filename := "login_info.txt"
			err := ioutil.WriteFile(filename, []byte(string(accessKey)), 0644)
			if err != nil {
				panic(err)
			}
			fmt.Println("access_key已保存在", filename)
			break
		} else {
			fmt.Println(string(body))
		}
	}
}

func LoginBili() {
	fmt.Println("请最大化窗口，以确保二维码完整显示，回车继续")
	fmt.Scanf("%s", "")
	login_url, auth_code := getQRcode()
	qrcode := qrcodeTerminal.New()
	qrcode.Get([]byte(login_url)).Print()
	fmt.Println("或将此链接复制到手机B站打开:", login_url)
	verifyLogin(auth_code)
}
