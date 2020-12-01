package util

import (
    "fmt"
    "os"
    "bytes"
    "reflect"
    "math/rand"
    "crypto/md5"
    "html/template"
    "encoding/hex"
    "time"
    "net"
    "unicode"
    "strings"
    "regexp"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "github.com/thinkeridea/go-extend/exnet"
)

func Unescaped (x string) interface{} { return template.HTML(x+"test") }


//通过反射判断某数值是否在一个slice中，
func IsExistItem(value interface{}, array interface{}) bool {
    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)
        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(value, s.Index(i).Interface()) {
                return true
            }
        }
    }
    return false
}
/*
    t := reflect.TypeOf(tmp)
    fmt.Println("tmp的类型是：", t)  // 打印：a的类型是： float64
*/
func FindType(i interface{}) {
    switch x := i.(type) {
    case int:
        fmt.Println(x, "is int")
    case string:
        fmt.Println(x, "is string")
    case nil:
        fmt.Println(x, "is nil")
    default:
        fmt.Println(x, "not type matched")
    }
}

func Substr(str string, start int, end int) string {
    rs := []rune(str)
    length := len(rs)

    if start < 0 || start > length {
        return ""
    }

    if end < 0 || end > length {
        return ""
    }
    return string(rs[start:end])
}

func GetMD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func Md5Check(content, encrypted string) bool {
    return strings.EqualFold(Md5Encode(content), encrypted)
}
func Md5Encode(data string) string {
    h := md5.New()
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}

func GenerateRangeNum(min, max int) int {
    rand.Seed(time.Now().Unix())
    randNum := rand.Intn(max - min) + min
    return randNum
}


func StrAssert(data interface {}) string {
    if str, ok := data.(string); ok {
        return str
    } else {
        return "nil"
    }
}

// transBooltostring
func Bts(b bool) string {
    if b {
        return "ok"
    } else {
        return "err"
    }
}

// check file if exist
func Exists(path string) bool {  
    _, err := os.Stat(path)    //os.Stat获取文件信息  
    if err != nil {  
        if os.IsExist(err) {  
            return true  
        }  
        return false  
    }  
    return true  
}  

func StrFilter(s string) string {
    mapping := func(r rune) rune {
        switch {
        case r >= 'A' && r <= 'Z': // 大写字母转小写
            return r + 32
        case r >= 'a' && r <= 'z': // 小写字母不处理
            return r
        case r >= '0' && r <= '9': // 数字不处理
            return r
        case unicode.Is(unicode.Han, r): // 汉字过滤掉
            return -1
        }
        return -1 // 过滤所有非字母、汉字的字符
    }
    out := strings.Map(mapping,s)
    return out
}

func parseDns(strDns string) (string, error) {

    ns, err := net.LookupHost(strDns)
    if err != nil {
        //fmt.Printf("error: %v, failed to parse %v\n", err, strDns)
        return strDns, err
    }
    //fmt.Printf("parse %v:\n", strDns)
    //for _, ip := range ns {
    //    fmt.Printf("%s\n", ip) 
    //}
    if len(ns) >= 1 {
        return ns[0], nil
    } else {
        return strDns, err
    }
}

func GetClientIp(r *http.Request) string{
    // var r *http.Request
    ip := exnet.ClientPublicIP(r)
    if ip == ""{
        ip = exnet.ClientIP(r)
    }
    return ip
}

func GetIp() string {
    var myurls = []string{"http://cip.cc/","https://ip.cn/","http://ip.cip.cc","http://myip.ipip.net"}
    ip := ""
    for _,myurl := range myurls {

        response, err := X509curlget(myurl)
        if err != nil {
            //fmt.Println(err.Error())
            return ip
        }
        defer response.Body.Close()

        data, err := ioutil.ReadAll(response.Body)
        if err != nil {
            return ip
        }
        //((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}
        re, _ := regexp.Compile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
        one := re.Find([]byte(data))
        ip = string(one)
        //fmt.Println(myurl,string(one))
        break;
    }
    return ip
}

// x509: certificate signed by unknown authority
// ref:https://zhangguanzhang.github.io/2019/07/07/golang-http/
// ref:https://studygolang.com/articles/12217
//func x509get(url string, headers map[string]string) *http.Response {
func X509get(url string) (*http.Response, error) {
    
    tr := &http.Transport{    //解决x509: certificate signed by unknown authority
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    
    client := &http.Client{
        Timeout:   60 * time.Second,
        Transport: tr,    //解决x509: certificate signed by unknown authority
    }
    req, err := http.NewRequest("GET", url, nil)

    //req.Header.Set("Accept-Encoding", "gzip, deflate, br")
    req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
    req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
    //req.Header.Set("Referer", baseUrl)
    req.Header.Set("Referer", "")
    req.Header.Set("X-Requested-With", "XMLHttpRequest")
    req.Header.Set("Connection", "keep-alive")

    //for k, v := range headers {
    //    req.Header.Add(k, v)
    //}

    if err != nil {
        fmt.Println(err.Error())
        return nil, err
    }

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err.Error())
        return nil, err
    }
    return resp, nil
}


func X509curlget(url string) (*http.Response, error) {
    tr := &http.Transport{    //解决x509: certificate signed by unknown authority
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{
        Timeout:   15 * time.Second,
        Transport: tr,    //解决x509: certificate signed by unknown authority
    }
    req, err := http.NewRequest("GET", url, nil)

    req.Header.Set("User-Agent", "curl/7.54.0")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Accept-Encoding", "")
    //for k, v := range headers {
    //    req.Header.Add(k, v)
    //}

    if err != nil {
        //log.Println(err.Error())
        return nil, err
    }

    resp, err := client.Do(req)
    if err != nil {
        //log.Println(err.Error())
        return nil, err
    }
    return resp, nil
}

func HttpPostJson(url string, jsonStr []byte) (int, string, error){
    //jsonStr =[]byte(`{ "username": "auto", "password": "auto123123" }`) 
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    if err != nil {
        return 0, "", err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return 0, "", err
    }
    defer resp.Body.Close()

    statuscode := resp.StatusCode
    //hea := resp.Header
    body, _ := ioutil.ReadAll(resp.Body)

    return statuscode, string(body), nil
    //fmt.Println(string(body))
    //fmt.Println(statuscode)
}
