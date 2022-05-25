package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	gjson "github.com/tidwall/gjson"
)

var AccessKey string

var Csrf string

var Cookies []*http.Cookie

func Login() {
	filename := "login_info.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil || len(data) == 0 {
		fmt.Println("未登录,请扫码登录")
		loginBili()
	} else {
		AccessKey = gjson.Parse(string(data)).Get("data.access_token").String()
		for _, c := range gjson.Parse(string(data)).Get("data.cookie_info.cookies").Array() {
			Cookies = append(Cookies, &http.Cookie{
				Name:  c.Get("name").String(),
				Value: c.Get("value").String(),
			})
			if c.Get("name").String() == "bili_jct" {
				Csrf = c.Get("value").String()
			}
		}
		l, name := is_login()
		if l {
			fmt.Println("登录成功：", name)
		} else {
			fmt.Println("登录失败，请重新扫码登录")
			loginBili()
		}
	}

}

func is_login() (bool, string) {
	api := "https://api.bilibili.com/x/web-interface/nav"
	client := http.Client{}
	req, _ := http.NewRequest("GET", api, nil)
	for _, c := range Cookies {
		req.AddCookie(c)
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
	data := make(map[string]string)
	data["local_id"] = "0"
	data["ts"] = fmt.Sprintf("%d", time.Now().Unix())
	signature(&data)
	data_string := strings.NewReader(map_to_string(data))
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
		panic("get_tv_qrcode_url_and_auth_code error")
	}
}

func verify_login(auth_code string) {
	api := "http://passport.bilibili.com/x/passport-tv-login/qrcode/poll"
	data := make(map[string]string)
	data["auth_code"] = auth_code
	data["local_id"] = "0"
	data["ts"] = fmt.Sprintf("%d", time.Now().Unix())
	signature(&data)
	data_string := strings.NewReader(map_to_string(data))
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
		AccessKey = gjson.Parse(string(body)).Get("data.access_token").String()
		if code == 0 {
			fmt.Println("登录成功")
			fmt.Println("access_key:", string(AccessKey))
			filename := "login_info.txt"
			err := ioutil.WriteFile(filename, []byte(string(AccessKey)), 0644)
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

var appkey = "4409e2ce8ffd12b8"
var appsec = "59b43e04ad6965f34319062b478f83dd"

func signature(params *map[string]string) {
	var keys []string
	(*params)["appkey"] = appkey
	for k := range *params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var query string
	for _, k := range keys {
		query += k + "=" + url.QueryEscape((*params)[k]) + "&"
	}
	query = query[:len(query)-1] + appsec
	hash := md5.New()
	hash.Write([]byte(query))
	(*params)["sign"] = hex.EncodeToString(hash.Sum(nil))
}

func map_to_string(params map[string]string) string {
	var query string
	for k, v := range params {
		query += k + "=" + v + "&"
	}
	query = query[:len(query)-1]
	return query
}

func loginBili() {
	fmt.Println("请最大化窗口，以确保二维码完整显示，回车继续")
	fmt.Scanf("%s", "")
	login_url, auth_code := get_tv_qrcode_url_and_auth_code()
	qrcode := qrcodeTerminal.New()
	qrcode.Get([]byte(login_url)).Print()
	fmt.Println("或将此链接复制到手机B站打开:", login_url)
	verify_login(auth_code)
}

func main() {
	loginBili()
	fmt.Scanf("%s", "")
}
