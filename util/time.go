package util

import (
    "fmt"
    "time"
    "strconv"
)


//@brief：耗时统计函数
// defer timeCost()()
func TimeCost() func() {
    start := time.Now()
    return func() {
        tc:=time.Since(start)
        //fmt.Printf("time cost = %v\n", tc)
        i64:=tc.Nanoseconds() / 1e6

        //logtc:=strconv.Itoa(int(i64))
        fmt.Println("time cost = ",int(i64),"ms")
    }
}

func Second2date(s int) string {
    if s < 60 {
        return "60s"
    }
    ret := strconv.Itoa(s)

    if (s >= 60) && (s < 3600) {
        m := s/60
        s1 := s%60

        str_m := strconv.Itoa(m)
        str_s := strconv.Itoa(s1)

        ret = str_m+"m"+str_s+"s"  
    }

    if (s >= 3600) && (s<86400) {
        h := s/3600
        sm := s%3600
        m := sm/60
        s1 := s%60
        str_h := strconv.Itoa(h)
        str_m := strconv.Itoa(m)
        str_s := strconv.Itoa(s1)

        ret = str_h+"h"+str_m+"m"+str_s+"s"
    }

    if s >= 86400 {
        d := s/86400
        sd := s%86400
        h := sd/3600
        sm := sd%3600
        m := sm/60
        s1 := s%60
        str_d := strconv.Itoa(d)
        str_h := strconv.Itoa(h)
        str_m := strconv.Itoa(m)
        str_s := strconv.Itoa(s1)

        ret = str_d+"d"+str_h+"h"+str_m+"m"+str_s+"s"
    }
    return ret
}



// Unix timestamp int 转换为 int64，再转换为时间字符串
func TimeInt2FormatStr(field int) string{
    // int 转换为 int64，再转换为时间字符串
    //int64, _ := strconv.ParseInt(field, 10, 64)  
    int64 := int64(field)  
    tm := time.Unix(int64, 0)
    //fmt.Println(tm.Format("2006-01-02 15:04:05"))
    ft := tm.Format("2006-01-02 15:04:05")
    return ft
}

//  Unix timestamp string 转换为 int64，再转换为时间字符串
func TimeString2FormatStr(field string) string{
    var cstZone = time.FixedZone("CST", 8*3600)       // 东八
    // int 转换为 int64，再转换为时间字符串
    int64, _ := strconv.ParseInt(field, 10, 64)  
    tm := time.Unix(int64, 0)
    //fmt.Println(tm.Format("2006-01-02 15:04:05"))
    ft := tm.In(cstZone).Format("2006-01-02 15:04:05")
    return ft
}

// 获取当前时间戳，并将时间戳转换为int
func GetCurrentTimestamp2int() int{
    var cstZone = time.FixedZone("CST", 8*3600)       // 东八

    timestamp := time.Now().In(cstZone).Unix()
    //intime := strconv.FormatInt(timestamp,10)
    intime := int(timestamp)

    return intime
}

// 获取当前时间戳，并将时间戳转换为string
func GetCurrentTimestamp2String() string{
    var cstZone = time.FixedZone("CST", 8*3600)       // 东八

    timestamp := time.Now().In(cstZone).Unix()
    intime := strconv.FormatInt(timestamp,10)

    return intime
}

func GetTimeString() string{
    var cstZone = time.FixedZone("CST", 8*3600)       // 东八

    tm := time.Now().In(cstZone)
    //intime := strconv.FormatInt(timestamp,10)
    ft := tm.In(cstZone).Format("2006-01-02 15:04:05")

    return ft
}

func String2Timestamp2Fmt(str string) string{
    //string 转 时间戳
    //stringTime := "20190702155040"
    stringTime := str
    loc, _ := time.LoadLocation("Local")
    the_time, err := time.ParseInLocation("20060102150405", stringTime, loc)
    if err != nil {
        return ""
    }
    unix_time := the_time.Unix() //1504082441
    //fmt.Println(unix_time)

    //时间戳转Time 再转 string
    timeNow := time.Unix(unix_time, 0) //2017-08-30 16:19:19 +0800 CST
    timeString := timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32
    //fmt.Println(timeString)
    return timeString
}

//时间"2006-01-02" to 时间戳
func FmtDateStr2UnixStr(str string) int64{
   loc, _ := time.LoadLocation("Asia/Shanghai")        //设置时区
   tt, _ := time.ParseInLocation("2006-01-02", str, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
   //fmt.Println(tt.Unix())
   return tt.Unix()*1000
}

//时间"09/23/2020" to 时间戳
func FmtDateStr2Unixint(str string) int{
   loc, _ := time.LoadLocation("Asia/Shanghai")        //设置时区
   tt, _ := time.ParseInLocation("01/02/2006", str, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
   //fmt.Println(tt.Unix())
   return int(tt.Unix())
}