package handler

import (
	"MedalHelper/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	gjson "github.com/tidwall/gjson"
)

var (
	accessKey string
	csrf      string
	cookies   []*http.Cookie
)

func Login() {
	filename := "login_info.json"
	if data, err := ioutil.ReadFile(filename); err != nil || len(data) == 0 {
		fmt.Println("未登录,请扫码登录")
		LoginBili()
	} else {
		accessKey = gjson.Parse(string(data)).Get("data.access_token").String()
		for _, c := range gjson.Parse(string(data)).Get("data.cookie_info.cookies").Array() {
			cookies = append(cookies, &http.Cookie{
				Name:  c.Get("name").String(),
				Value: c.Get("value").String(),
			})
			if c.Get("name").String() == "bili_jct" {
				csrf = c.Get("value").String()
			}
		}
		l, name := isLogin()
		if l {
			fmt.Println("登录成功：", name)
		} else {
			fmt.Println("登录失败，请重新扫码登录")
			LoginBili()
		}
	}
}

func isLogin() (bool, string) {
	api := "https://api.bilibili.com/x/web-interface/nav"
	client := http.Client{}
	req, _ := http.NewRequest("GET", api, nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := gjson.ParseBytes(body)
	return data.Get("code").Int() == 0, data.Get("data.uname").String()
}

func get_tv_qrcode_url_and_auth_code() (string, string) {
	api := "http://passport.bilibili.com/x/passport-tv-login/qrcode/auth_code"
	data := map[string]string{
		"local_id": "0",
		"tx":       fmt.Sprintf("%d", time.Now().Unix()),
	}
	util.Signature(&data)
	data_string := strings.NewReader(util.Map2string(data))
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
		panic("get_tv_qrcode_url_and_auth_code error")
	}
}

func verifyLogin(auth_code string) {
	api := "http://passport.bilibili.com/x/passport-tv-login/qrcode/poll"
	data := map[string]string{
		"auth_code": auth_code,
		"local_id":  "0",
		"ts":        fmt.Sprintf("%d", time.Now().Unix()),
	}
	util.Signature(&data)
	data_string := strings.NewReader(util.Map2string(data))
	client := http.Client{}
	req, _ := http.NewRequest("POST", api, data_string)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for {
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
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
			time.Sleep(time.Second * 3)
		}
	}
}

func LoginBili() {
	fmt.Println("请最大化窗口，以确保二维码完整显示，回车继续")
	fmt.Scanf("%s", "")
	login_url, auth_code := get_tv_qrcode_url_and_auth_code()
	qrcode := qrcodeTerminal.New()
	qrcode.Get([]byte(login_url)).Print()
	fmt.Println("或将此链接复制到手机B站打开:", login_url)
	verifyLogin(auth_code)
}
