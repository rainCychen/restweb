package controller

import (
	"controller/test"
	"net"
	"net/http"
	"strconv"

	"time"

	"encoding/json"

	"io"

	"fmt"

	"github.com/emicklei/go-restful"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
)

type Config struct {
	Port     int
	DbName   string
	RedisSvr string
	RedisPwd string
	RedisDb  int
}
type Controller struct {
	Cfg        Config
	Container  *restful.Container
	Container2 *restful.Container
}

const (
	StaicDir  string = "src/static/"
	StaicCss  string = "src/static/css/"
	StaicJs   string = "src/static/js/"
	StaicHtml string = "src/static/html/"
)

func New(f Config) *Controller {
	return &Controller{Cfg: f}
}
func newPool(server, password string, ma int, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		MaxActive:   ma,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", server, 3*time.Second, 5*time.Second, 5*time.Second)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

var redisPool *redis.Pool

func (c *Controller) StartUp() error {
	redisPool = newPool(c.Cfg.RedisSvr, c.Cfg.RedisPwd, 2, c.Cfg.RedisDb)
	c.Container = restful.NewContainer()
	restful.PrettyPrintResponses = true
	p := &test.ProductResource{}
	web := test.New(StaicHtml, redisPool)
	//允许 被跨域请求到
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization", "Range", "Origin", "AD-Expire"},
		ExposeHeaders:  []string{"Content-Length", "content-type"},
		AllowedMethods: []string{"POST", "PUT", "DELETE", "GET"},
		CookiesAllowed: false,
		Container:      c.Container,
	}
	c.Container.Filter(cors.Filter)

	web.InitRoute(c.Container)
	p.InitRoute(c.Container)

	ws1 := new(restful.WebService)
	ws1.Path("/v3/healthz")
	ws1.ApiVersion("v3")
	ws1.Produces(restful.MIME_JSON)
	ws1.Consumes(restful.MIME_JSON)
	ws1.Route(ws1.GET("/").To(handleHealthz))
	ws1.Route(ws1.GET("/send").To(handleSend))
	c.Container.Add(ws1)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(StaicDir)))) //设置css、js存放的路径
	svr := http.Server{Handler: c.Container, Addr: net.JoinHostPort("0.0.0.0", strconv.Itoa(c.Cfg.Port))}
	return svr.ListenAndServe()
}
func handleHealthz(req *restful.Request, rsp *restful.Response) {
	//rsp.ResponseWriter.Write([]byte("ok"))
	conn := redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	m := make(map[string]interface{})
	if err == nil {
		m["Flag"] = 100
		m["redis"] = "ok"
	} else {
		m["Flag"] = 500
		m["redis"] = "error:" + err.Error()
	}
	rsp.WriteEntity(m)
}

type Test struct {
	a int
	b string
	c bool
}

func handleSend(req *restful.Request, rsp *restful.Response) {
	//rsp.ResponseWriter.Write([]byte("ok"))
	conn := redisPool.Get()
	defer conn.Close()
	//后端绘制队列
	qm := &Test{
		a: 1,
		b: "redis",
		c: true,
	}
	buf, _ := json.Marshal(qm)
	_, err := conn.Do("LPUSH", "name", "111", 234, false)
	if err != nil {
		glog.Warningln("LPUSH fail", err)
	}
	//	err = conn.Send("LPUSH", "test", buf)
	//	if err != nil {
	//		glog.Warningln("LPUSH send fail", err)
	//	}
	glog.Infof("消息%v %v\n", qm, string(buf))
	//	rsp.WriteEntity(qm)
	io.WriteString(rsp, fmt.Sprintf("a%vb%vc%v", qm.a, qm.b, qm.c))
}
