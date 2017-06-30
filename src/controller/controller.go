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
	cfg       Config
	container *restful.Container
}

func New(f Config) *Controller {
	return &Controller{cfg: f}
}
func (c *Controller) StartUp() error {
	c.container = restful.NewContainer()
	restful.PrettyPrintResponses = false
	web := &test.Web{}
	web.InitRoute(c.container)
	p := &test.ProductResource{}
	p.InitRoute(c.container)
	svr := http.Server{Handler: c.container, Addr: net.JoinHostPort("0.0.0.0", strconv.Itoa(c.cfg.Port))}
	return svr.ListenAndServe()
}
