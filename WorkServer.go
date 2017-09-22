package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
	"encoding/json"

	"time"
)

type WorkServer struct {
	serverData RunningData
	masterdata RunningData
	running    bool
	zkClient   *zk.Conn

	path       string
}
var zkServersAdd =[]string{"192.168.0.21:2182"}

func NewWorkServer(rData RunningData) (*WorkServer) {
	ws := new(WorkServer)
	ws.serverData = rData
	ws.running = false
	var err error
	ws.zkClient,err = Connect()
	if err!=nil {
		fmt.Println(err.Error())
		return nil
	}
	ws.path = "/master"
	return ws
}

func (w *WorkServer) start() {
	if w.running {
		fmt.Println("server has startup...")
		return
	}
	w.running = true
	w.takerMaster()
}

func Connect()(conn *zk.Conn,err error){
	conn,_,err=zk.Connect(zkServersAdd,time.Second * 3)
	return
}

func (w *WorkServer) stop() {
	if !w.running {
		fmt.Println("server has stoped...")
		return
	}
	w.running = false
	w.releaseMaster()
}

func (w *WorkServer) takerMaster() {
	if !w.running {
		return
	}

	jd, _ := json.Marshal(w.serverData)
	_, err := w.zkClient.Create(w.path, jd, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	if err == nil { // 创建成功表示选举master成功
		w.masterdata = w.serverData
		fmt.Println(time.Now(),w.serverData.Cid,"获得master")
	} else { // 创建失败表示选举master失败
		//读取节点
		path, _, event, e := w.zkClient.GetW(w.path)
		if e == nil {
			json.Unmarshal(path, w.masterdata)
			go w.event(event)
		} else {
			w.takerMaster()
		}
	}

}

func (w *WorkServer) releaseMaster() {
	if w.cleckMaster() {
		w.zkClient.Delete(w.path, 0)
	}
	w.running=false
}

func (w *WorkServer) cleckMaster() bool {

	path, _, e :=w.zkClient.Get(w.path)
	if e==nil {
		var eventData RunningData
		json.Unmarshal(path, eventData)

		if w.serverData.Cid==eventData.Cid {
			return true
		}else {
			return  false
			}

	}else{
		return false
	}
}

func (w *WorkServer) event(event <-chan zk.Event) {
	select  {
		case e:=<-event:
			if e.Type==zk.EventNodeDeleted{
				w.takerMaster()
			}else if e.Type==zk.EventNodeDataChanged {


			}

	}
}
