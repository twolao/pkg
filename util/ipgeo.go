package util

import (
	"github.com/ipipdotnet/ipdb-go"
)

func Ip2addr(ip string) string {
	db, err := ipdb.NewCity("./conf/ipipfree.ipdb")
	if err != nil {
		return ""
	}
	//db.Reload("/path/to/city.ipv4.ipdb") // 更新 ipdb 文件后可调用 Reload 方法重新加载内容

	ret, err := db.FindMap(ip, "CN") // return map[string]string
	if err!=nil {
		return ""
	}
	out := ""
	if ret["country_name"] == ret["region_name"] {
		out = ret["country_name"]
		if ret["city_name"] != ret["region_name"] {
			out += ret["city_name"]	
		}
		
	} else {
		out = ret["country_name"] + ret["region_name"]
		if ret["city_name"] != ret["region_name"] {
			out += ret["city_name"]	
		}
	}

	return out
}