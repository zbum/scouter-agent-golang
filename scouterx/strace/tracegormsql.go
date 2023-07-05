package strace

import (
	"context"
	"fmt"
	"github.com/zbum/scouter-agent-golang/scouterx/strace/tctxmanager"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	callbackName = "gorm:%s"
	registerName = "scouter:%s"
)

type GormDbPlugin struct {
}

func name(pattern, value string) string {
	return fmt.Sprintf(pattern, value)
}
func (s GormDbPlugin) Name() string {
	return "scouter-plugin"
}

func (s GormDbPlugin) Initialize(db *gorm.DB) error {
	s.registerCreate(db)
	s.registerQuery(db)
	s.registerUpdate(db)
	s.registerDelete(db)
	return nil
}

func (s GormDbPlugin) registerCreate(db *gorm.DB) {
	callback := name(callbackName, "create")

	_ = db.Callback().Create().Before(callback).
		Register("before_create", before("INSERT"))

	_ = db.Callback().Create().After(callback).
		Register("after_create", after("INSERT"))
}
func (s GormDbPlugin) registerQuery(db *gorm.DB) {
	callback := name(callbackName, "query")

	_ = db.Callback().Query().Before(callback).
		Register("before_query", before("SELECT"))

	_ = db.Callback().Query().After(callback).
		Register("after_query", after("SELECT"))
}

func (s GormDbPlugin) registerUpdate(db *gorm.DB) {
	callback := name(callbackName, "update")

	_ = db.Callback().Update().Before(callback).
		Register("before_update", before("UPDATE"))

	_ = db.Callback().Update().After(callback).
		Register("after_update", after("UPDATE"))
}

func (s GormDbPlugin) registerDelete(db *gorm.DB) {
	callback := name(callbackName, "delete")

	_ = db.Callback().Delete().Before(callback).
		Register("before_delete", before("DELETE"))

	_ = db.Callback().Delete().After(callback).
		Register("after_delete", after("DELETE"))
}

var before = func(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		tctx := tctxmanager.GetTraceContext(db.Statement.Context)
		db.Statement.Context = context.WithValue(db.Statement.Context, fmt.Sprintf("%d-%s", tctx.Txid, "stepStartTime"), time.Now().UnixMilli())
	}
}

var after = func(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		tctx := tctxmanager.GetTraceContext(db.Statement.Context)
		starTime := db.Statement.Context.Value(fmt.Sprintf("%d-%s", tctx.Txid, "stepStartTime")).(int64)
		tctx.SqlCount++
		tctx.SqlTime += int32(time.Now().UnixMilli() - starTime)
		AddMessageStep(db.Statement.Context,
			fmt.Sprintf("STM> %s", db.Statement.SQL.String()))
		AddMessageStep(db.Statement.Context,
			fmt.Sprintf("[%s]", joinSlice(db.Statement.Vars, ",")))
	}
}

func joinSlice(slice []interface{}, sep string) (result string) {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprintf("%v", v)
	}

	result = strings.Join(strSlice, sep)
	return
}
