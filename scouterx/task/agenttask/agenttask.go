package agenttask

import (
	"github.com/zbum/scouter-agent-golang/scouterx/common/logger"
	"github.com/zbum/scouter-agent-golang/scouterx/common/netdata"
	"github.com/zbum/scouter-agent-golang/scouterx/conf"
	"github.com/zbum/scouter-agent-golang/scouterx/netio"
	"time"
)

func Start() {
	go start()
}

func start() {
	for {
		time.Sleep(2000 * time.Millisecond)
		SendObjPack()
	}
}

func SendObjPack() {
	ac := conf.GetInstance()
	objName := ac.ObjName
	objHash := netio.SendObjName(objName)

	objPack := netdata.NewObjectPack()
	objPack.ObjName = objName
	objPack.ObjHash = objHash
	objPack.ObjType = ac.ObjType
	objPack.Version = "0.0.0"

	if ac.TraceObjSend {
		logger.Info.Printf("[scouter] SendObjPack: %s, %d, %s", objName, objHash, ac.ObjType)
	}
	netio.SendPackDirect(objPack)
}
