package utils

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
)

type IpQuery struct {
	Query      string `json:"query"`      //查询IP
	Status     string `json:"status"`     //查询状态
	Country    string `json:"country"`    //国家
	RegionName string `json:"regionName"` //省
	City       string `json:"city"`       //城市
	Org        string `json:"org"`        //组织
}

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
)

// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
// 来源: https://github.com/polaris1119/goutils/blob/master/ip.go
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

// Ip2long 将 IPv4 字符串形式转为 uint32
// 来源: https://github.com/polaris1119/goutils/blob/master/ip.go
func Ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

// IpGeolocation 返回IP归属地信息
// 调用 ip-api.com的api服务
func IpGeolocation(ip string) string {
	resp, err := http.Get("http://ip-api.com/json/" + ip + "?lang=zh-CN&fields=status,message,country,regionName,city,org,query")
	if err != nil {
		return "未知的地址"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "未知的地址"
	}
	var q = new(IpQuery)
	if err := json.Unmarshal(body, q); err != nil {
		return "未知的地址"
	}
	if q.Status == "success" {
		return q.Country + " " + q.RegionName + " " + q.City + " (" + q.Org + ")"
	} else {
		return "未知的地址"
	}
	return "未知的地址"
}
