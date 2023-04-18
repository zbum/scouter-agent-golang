package scouterx

import (
	"github.com/zbum/scouter-agent-golang/scouterx/common"
	"github.com/zbum/scouter-agent-golang/scouterx/conf"
	"github.com/zbum/scouter-agent-golang/scouterx/netio"
	"github.com/zbum/scouter-agent-golang/scouterx/strace"
	"github.com/zbum/scouter-agent-golang/scouterx/task/agenttask"
	"github.com/zbum/scouter-agent-golang/scouterx/task/countertask"
	"sync"
)

type objNameChangeObserver struct{}

var agentOnce sync.Once
var observer objNameChangeObserver

func (objNameChangeObserver) Run() {
	ac := conf.GetInstance()
	netio.SendObjName(ac.ObjName)
}

func init() {
	Init()
}

func Init() {
	agentOnce.Do(func() {
		common.RegisterSendPanicF(netio.SendPanic)
		conf.GetInstance()
		strace.StartTracingMode()
		agenttask.SendObjPack()
		agenttask.Start()
		conf.AddObjChangedObserver(observer)
		countertask.Start()
	})
}
