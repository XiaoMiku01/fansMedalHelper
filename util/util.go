package util

import (
	"fmt"
	"net/url"
	"time"
)

// TODO: optimize this part
// Map2String can transfer a string-string map into a raw string
func Map2string(params map[string]string) string {
	var query string
	for k, v := range params {
		query += k + "=" + v + "&"
	}
	query = query[:len(query)-1]
	return query
}

// Map2String can transfer a string-string map into url value struct
func Map2Params(params map[string]string) url.Values {
	value := url.Values{}
	for key, param := range params {
		value[key] = []string{param}
	}
	return value
}

// GetTimestamp can obtain current ts
func GetTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}
