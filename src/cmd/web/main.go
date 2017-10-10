package main

import (
	"controller"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/golang/glog"
)

var (
	port     = flag.Int("port", 8800, "port")
	dbName   = flag.String("db_name", "yun_wx2", "")
	RedisSvr = flag.String("qrsvr", "127.0.0.1:6379", "redis queue svr addr")
	RedisPwd = flag.String("qrpwd", "cz125678", "redis queue pwd ")
	RedisDb  = flag.Int("qrdb", 0, "redis queue db ")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	go func() {
		if err := http.ListenAndServe(":", nil); err != nil {
			glog.Errorln("pprof", err)
		}
	}()
	config := controller.Config{
		Port:     *port,
		DbName:   *dbName,
		RedisSvr: *RedisSvr,
		RedisPwd: *RedisPwd,
		RedisDb:  *RedisDb,
	}
	c := controller.New(config)
	err := c.StartUp()
	glog.Warningln("start fail", err)
}
