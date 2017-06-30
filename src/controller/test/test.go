package test

import (
	//"html/template"
	"github.com/emicklei/go-restful"
	//"github.com/golang/glog"
	//"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Web struct{}

func (w *Web) InitRoute(Container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1")
	ws.ApiVersion("v1")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(w.handleTest))
	ws.Route(ws.GET("/index").To(w.handleIndex))
	ws.Route(ws.GET("/login").To(w.handleForm))
	ws.Route(ws.POST("/login").Consumes("application/x-www-form-urlencoded").To(w.handleLogin))
	ws.Route(ws.GET("/user").To(w.handleForm))
	ws.Route(ws.POST("/login").Consumes("application/x-www-form-urlencoded").To(w.handleUser))
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

func (w *Web) handleTest(req *restful.Request, rsp *restful.Response) {
	info := Info{
		Id:   1,
		Name: "小明",
		Msg:  "hello world",
	}
	//glog.Infof("dafdsfa%v\n", info)
	fmt.Printf("%v\n", info)
	m := make(map[string]interface{})
	m["list"] = info
	io.WriteString(rsp, "小明")
}
func (w *Web) handleIndex(req *restful.Request, rsp *restful.Response) {
	info := Info{
		Id:   1,
		Name: "小明",
		Msg:  "hello world",
	}
	rsp.WriteEntity(info)
}
func (w *Web) handleLogin(req *restful.Request, rsp *restful.Response) {
	err := req.Request.ParseForm()
	//	p, err := url.QueryUnescape(req.URL.RawQuery)
	if err != nil {
		rsp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	p := new(Profile)
	p.Name = req.Request.PostFormValue("Name")
	p.Age, err = strconv.Atoi(req.Request.PostFormValue("Age"))
	//	err = json.Unmarshal(req.Request.PostForm.([]byte), &p)
	//	err = decoder.Decode(p, req.Request.PostForm)
	if err != nil {
		rsp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	//	fmt.Printf("%T(%[1]v)\n", req.Request.PostFormValue("Name"))
	//	rsp.WriteEntity(req.Request.PostFormValue("Age"))
	io.WriteString(rsp.ResponseWriter, fmt.Sprintf("<html><body>Name=%s, Age=%d</body></html>", p.Name, p.Age))
}
func (w *Web) handleUser(req *restful.Request, rsp *restful.Response) {
	p, err := url.QueryUnescape(url.URL.RawQuery)
	if err != nil {
		rsp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	p := new(Profile)
	p.Name = req.Request.PostFormValue("Name")
	p.Age, err = strconv.Atoi(req.Request.PostFormValue("Age"))
	//	err = json.Unmarshal(req.Request.PostForm.([]byte), &p)
	//	err = decoder.Decode(p, req.Request.PostForm)
	if err != nil {
		rsp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	//	fmt.Printf("%T(%[1]v)\n", req.Request.PostFormValue("Name"))
	//	rsp.WriteEntity(req.Request.PostFormValue("Age"))
	io.WriteString(rsp.ResponseWriter, fmt.Sprintf("<html><body>Name=%s, Age=%d</body></html>", p.Name, p.Age))
}
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