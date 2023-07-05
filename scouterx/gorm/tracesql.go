package gorm

import (
	"fmt"
	"github.com/zbum/scouter-agent-golang/scouterx/strace"
	"gorm.io/gorm"
)

const (
	callbackName = "gorm:%s"
	registerName = "scouter:%s"
)

type DbPlugin struct {
}

func name(pattern, value string) string {
	return fmt.Sprintf(pattern, value)
}
func (s DbPlugin) Name() string {
	return "scouter-plugin"
}

func (s DbPlugin) Initialize(db *gorm.DB) error {
	s.registerCreate(db)
	s.registerQuery(db)
	s.registerUpdate(db)
	s.registerDelete(db)
	return nil
}

func (s DbPlugin) registerCreate(db *gorm.DB) {
	callback := name(callbackName, "create")

	_ = db.Callback().Create().After(callback).
		Register("after_create", after("INSERT"))
}
func (s DbPlugin) registerQuery(db *gorm.DB) {
	callback := name(callbackName, "query")

	_ = db.Callback().Query().After(callback).
		Register("after_query", after("SELECT"))
}

func (s DbPlugin) registerUpdate(db *gorm.DB) {
	callback := name(callbackName, "update")

	_ = db.Callback().Update().After(callback).
		Register("after_update", after("UPDATE"))
}

func (s DbPlugin) registerDelete(db *gorm.DB) {
	callback := name(callbackName, "delete")

	_ = db.Callback().Delete().After(callback).
		Register("after_delete", after("DELETE"))
}

var after = func(operation string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		strace.AddMessageStep(db.Statement.Context,
			fmt.Sprintf("%s", db.Statement.SQL.String()))
		strace.AddMessageStep(db.Statement.Context,
			fmt.Sprintf("%v", db.Statement.Vars))
	}
}
