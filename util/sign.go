package util

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
)

const (
	AppKey = "4409e2ce8ffd12b8"
	AppSec = "59b43e04ad6965f34319062b478f83dd"
)

func Signature(params *map[string]string) {
	var keys []string
	(*params)["appkey"] = AppKey
	for k := range *params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var query string
	for _, k := range keys {
		query += k + "=" + url.QueryEscape((*params)[k]) + "&"
	}
	query = query[:len(query)-1] + AppSec
	hash := md5.New()
	hash.Write([]byte(query))
	(*params)["sign"] = hex.EncodeToString(hash.Sum(nil))
}
