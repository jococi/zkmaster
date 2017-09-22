package main

import (
	"fmt"
	"github.com/koangel/grapeTimer"
)

var MapServer map[int] *WorkServer

func main() {
	grapeTimer.InitGrapeScheduler(1000,true)
	grapeTimer.NewTickerLoop(5000,-1,sssTimer,nil)

	MapServer=make(map[int] *WorkServer)
	for i:=0;i<10 ;i++  {
		ws:=NewWorkServer(RunningData{int64(i),"no1"})
		if ws!=nil {
			ws.start()
			MapServer[i]=ws
		}

	}




	signalCH := InitSignal()
	HandleSignal(signalCH)
	// exit
	fmt.Println("zk stop")
}

func sssTimer(timerId int, args interface{}) {
	for _, s := range MapServer {
		if s.cleckMaster() {
			s.releaseMaster()
			s.start()
		}
	}
}
