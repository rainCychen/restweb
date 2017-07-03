package controller

import (
	"controller/test"
	"net"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
)

type Config struct {
	Port   int
	DbName string
}
type Controller struct {
	Cfg       Config
	Container *restful.Container
}

func New(f Config) *Controller {
	return &Controller{Cfg: f}
}
func (c *Controller) StartUp() error {
	c.Container = restful.NewContainer()
	restful.PrettyPrintResponses = false
	web := &test.Web{}
	web.InitRoute(c.Container)
	p := &test.ProductResource{}
	p.InitRoute(c.Container)
	svr := http.Server{Handler: c.Container, Addr: net.JoinHostPort("0.0.0.0", strconv.Itoa(c.Cfg.Port))}
	return svr.ListenAndServe()
}
