package strace

import (
	"github.com/zbum/scouter-agent-golang/scouterx/common/logger"
	"github.com/zbum/scouter-agent-golang/scouterx/common/netdata"
	"github.com/zbum/scouter-agent-golang/scouterx/common/util"
	"github.com/zbum/scouter-agent-golang/scouterx/conf"
	"github.com/zbum/scouter-agent-golang/scouterx/netio"
	"sync"
	"testing"
	"time"
)

func TestStartTracingMode(t *testing.T) {
	ac := conf.GetInstance()
	ac.SetTrace(true)

	logger.Error.Println("error log test")
	logger.Error.Printf("error log test %s\n", "(testing)")
	wg := sync.WaitGroup{}
	wg.Add(1)

	objPack := netdata.NewObjectPack()
	objPack.ObjName = "node-testcase-start"
	objPack.ObjHash = util.HashString(objPack.ObjName)
	objPack.ObjType = "java"
	netio.SendPackDirect(objPack)
	ac.ObjHash = objPack.ObjHash

	go func() {
		for true {
			time.Sleep(3000 * time.Millisecond)
			netio.SendPackDirect(objPack)
		}
	}()
	StartTracingMode()

	wg.Wait()
}
