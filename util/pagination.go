package util

import (
	"fmt"
    "math"
    "net/http"
    "net/url"
    "strconv"
    "strings"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/oceanchang/pkg/setting"
)

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
    if page > 0 {
        result = (page - 1) * setting.PageSize
    }

    return result
}

//为了mongodb请求函数特殊处理
func GetPagem(c *gin.Context) int {
    result := 0
    page, _ := com.StrTo(c.Query("page")).Int()
    if page > 0 {
        result = page - 1
    }

    return result
}


//Pagination 分页器
type Pagination struct {
    Request *http.Request
    Total   int
    Pernum  int
}

//NewPagination 新建分页器
func NewPagination(req *http.Request, total int, pernum int) *Pagination {
    return &Pagination{
        Request: req,
        Total:   total,
        Pernum:  pernum,
    }
}

//Pages 渲染生成html分页标签
func (p *Pagination) Pages() string {
    queryParams := p.Request.URL.Query()
    //从当前请求中获取page
    page := queryParams.Get("page")
    if page == "" {
        page = "1"
    }
    //将页码转换成整型，以便计算
    pagenum, _ := strconv.Atoi(page)
    if pagenum == 0 {
        return ""
    }

    //计算总页数
    var totalPageNum = int(math.Ceil(float64(p.Total) / float64(p.Pernum)))

    //首页链接
    var firstLink string
    //上一页链接
    var prevLink string
    //下一页链接
    var nextLink string
    //末页链接
    var lastLink string
    //中间页码链接
    var pageLinks []string

    //首页和上一页链接
    if pagenum > 1 {
        firstLink = fmt.Sprintf(`<li><a class="btn btn-default btn-sm" href="%s">首页</a></li>`, p.pageURL("1"))
        prevLink = fmt.Sprintf(`<li><a class="btn btn-default btn-sm" href="%s">上一页</a></li>`, p.pageURL(strconv.Itoa(pagenum-1)))
    } else {
        firstLink = `<li class="disabled"><a class="btn btn-default btn-sm" href="#">首页</a></li>` 
        prevLink = `<li class="disabled"><a class="btn btn-default btn-sm" href="#">上一页</a></li>`
    }

    //末页和下一页
    if pagenum < totalPageNum {
        lastLink = fmt.Sprintf(`<li><a class="btn btn-default btn-sm" href="%s">末页</a></li>`, p.pageURL(strconv.Itoa(totalPageNum)))
        nextLink = fmt.Sprintf(`<li><a class="btn btn-default btn-sm" href="%s">下一页</a></li>`, p.pageURL(strconv.Itoa(pagenum+1)))
    } else {
        lastLink = `<li class="disabled"><a class="btn btn-default btn-sm" href="#">末页</a></li>`
        nextLink = `<li class="disabled"><a class="btn btn-default btn-sm" href="#">下一页</a></li>`
    }

    //生成中间页码链接
    pageLinks = make([]string, 0, 10)
    startPos := pagenum - 3
    endPos := pagenum + 3
    if startPos < 1 {
        endPos = endPos + int(math.Abs(float64(startPos))) + 1
        startPos = 1
    }
    if endPos > totalPageNum {
        endPos = totalPageNum
    }
    for i := startPos; i <= endPos; i++ {
        var s string
        if i == pagenum {
            s = fmt.Sprintf(`<li class="active"><a class="btn btn-default btn-sm" href="%s">%d</a></li>`, p.pageURL(strconv.Itoa(i)), i)
        } else {
            s = fmt.Sprintf(`<li><a class="btn btn-default btn-sm" href="%s">%d</a></li>`, p.pageURL(strconv.Itoa(i)), i)
        }
        pageLinks = append(pageLinks, s)
    }

    return fmt.Sprintf(`<ul class="pagination">%s%s%s%s%s</ul>`, firstLink, prevLink, strings.Join(pageLinks, ""), nextLink, lastLink)
}


//pageURL 生成分页url
func (p *Pagination) pageURL(page string) string {
    //基于当前url新建一个url对象
    u, _ := url.Parse(p.Request.URL.String())
    q := u.Query()
    q.Set("page", page)
    u.RawQuery = q.Encode()
    return u.String()
}