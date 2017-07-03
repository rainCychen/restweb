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
func (c *Controller) StartUp() error {
	c.Container = restful.NewContainer()
	restful.PrettyPrintResponses = true
	//允许 被跨域请求到
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization", "Range", "Origin", "AD-Expire"},
		ExposeHeaders:  []string{"Content-Length", "content-type"},
		AllowedMethods: []string{"POST", "PUT", "DELETE", "GET"},
		CookiesAllowed: false,
		Container:      c.Container,
	}
	c.Container.Filter(cors.Filter)
	test.New(StaicHtml)
	web := &test.Web{}
	web.InitRoute(c.Container)
	p := &test.ProductResource{}
	p.InitRoute(c.Container)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(StaicDir)))) //设置css、js存放的路径
	svr := http.Server{Handler: c.Container, Addr: net.JoinHostPort("0.0.0.0", strconv.Itoa(c.Cfg.Port))}
	return svr.ListenAndServe()
}
