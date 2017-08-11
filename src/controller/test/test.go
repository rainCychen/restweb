package test

import (
	"encoding/json"
	"html/template"

	"github.com/emicklei/go-restful"
	//	"github.com/golang/glog"
	//"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	//	"text/template"
	"github.com/garyburd/redigo/redis"
)

type Web struct {
	redisPool *redis.Pool
	StaicHtml string
}

func New(route string, redisPool *redis.Pool) *Web {
	return &Web{
		redisPool: redisPool,
		StaicHtml: route,
	}

}
func (w *Web) InitRoute(Container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1")
	//	ws.Produces(restful.MIME_JSON)
	//	ws.Consumes(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(w.handleTest))
	ws.Route(ws.GET("/index").To(w.handleIndex))
	ws.Route(ws.GET("/login").To(w.handleForm))
	ws.Route(ws.POST("/login").Consumes("application/x-www-form-urlencoded").To(w.handleLogin))
	ws.Route(ws.GET("/user").To(w.handleUserForm))
	ws.Route(ws.POST("/user").Consumes("application/x-www-form-urlencoded").To(w.handleUser))
	ws.Route(ws.GET("/html").To(w.handleHtml))
	ws.Route(ws.GET("/redis").To(w.handleSync))
	Container.Add(ws)

}

type Info struct {
	Id   int
	Name string
	Msg  string
}
type Profile struct {
	Name string
	Age  int
}

/*
简单的信息打印（io.WriteString(rsp, "小明")）
*/
func (w *Web) handleTest(req *restful.Request, rsp *restful.Response) {
	info := Info{
		Id:   1,
		Name: "小明",
		Msg:  "hello world",
	}
	//glog.Infof("dafdsfa%v\n", info)
	fmt.Printf("%v\n", info)
	m := make(map[string]interface{})
	m["list"] = info.Name
	//	io.WriteString(rsp, info.Msg)
	byte, _ := json.Marshal(info)
	io.WriteString(rsp, string(byte))
}

/*
简单的信息打印（rsp.WriteEntity(info)）
*/
func (w *Web) handleIndex(req *restful.Request, rsp *restful.Response) {
	info := Info{
		Id:   1,
		Name: "小明",
		Msg:  "hello world",
	}
	rsp.WriteEntity(info)
}

/*
表单数据
*/
func (w *Web) handleLogin(req *restful.Request, rsp *restful.Response) {
	err := req.Request.ParseForm()
	if err != nil {
		rsp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	p := new(Profile)
	p.Name = req.Request.PostFormValue("Name")
	p.Age, err = strconv.Atoi(req.Request.PostFormValue("Age"))
	if err != nil {
		rsp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	io.WriteString(rsp.ResponseWriter, fmt.Sprintf("<html><body>Name=%s, Age=%d</body></html>", p.Name, p.Age))
}

/*
表单
*/
func (w *Web) handleForm(req *restful.Request, rsp *restful.Response) {
	io.WriteString(rsp,
		`<html>
		<body>
		<h1>Enter Profile</h1>
		<form method="post">
		    <label>Name:</label>
			<input type="text" name="Name"/>
			<label>Age:</label>
		    <input type="text" name="Age"/>
			<input type="Submit" />
		</form>
		</body>
		</html>`)
}

/*
使用模版
*/

func (w *Web) handleHtml(req *restful.Request, rsp *restful.Response) {
	t, err := template.ParseFiles(w.StaicHtml + "hello.html")
	if err != nil {
		log.Fatalf("Template gave: %s", err)
	}
	t.Execute(rsp.ResponseWriter, nil)
}
func (w *Web) handleUserForm(req *restful.Request, rsp *restful.Response) {
	t, err := template.ParseFiles(w.StaicHtml + "user.html")
	if err != nil {
		log.Fatalf("Template gave: %s", err)
	}
	t.Execute(rsp.ResponseWriter, nil)
}

/*
获取表单get中的数据
*/
func (w *Web) handleUser(req *restful.Request, rsp *restful.Response) {
	id, err := url.QueryUnescape(req.Request.URL.RawQuery)
	if err != nil {
		rsp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	err = req.Request.ParseForm()
	if err != nil {
		rsp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	p := new(Profile)
	p.Name = req.QueryParameter("Name")
	p.Age, err = strconv.Atoi(req.QueryParameter("Age")) //字符串转整形
	fmt.Printf("%T(%[1]v)", id)
	io.WriteString(rsp, fmt.Sprintf("GET:%T(%[1]v)Name%vAge%v", id, p.Name, p.Age))
}

/* 用redis发送消息 */
func (w *Web) handleSync(req *restful.Request, rsp *restful.Response) {
	a, b, c := "test", 1, true
	err := w.Sync(a, b, c)
	if err != nil {
		io.WriteString(rsp, fmt.Sprintf("redis error: %v", err))
	} else {
		io.WriteString(rsp, "redis")
	}

}
