package test

import (
	"encoding/json"

	"github.com/golang/glog"
)

type Queue struct {
	A string `json:"a"`
	B int    `json:"b"`
	C bool   `json:"c"`
}

/* 用redis发送消息 */
func (w *Web) Sync(a string, b int, c bool) error {
	conn := w.redisPool.Get()
	defer conn.Close()
	//后端绘制队列
	qm := Queue{
		A: a,
		B: b,
		C: c,
	}
	buf, _ := json.Marshal(qm)
	glog.Infoln(buf, qm)
	_, err := conn.Do("LPUSH", "test", buf)
	if err != nil {
		glog.Warningln("LPUSH fail", err)
		return err
	}
	return nil
}
