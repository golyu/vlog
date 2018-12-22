package vlog

// BindGormLogger 绑定gorm的日志
//noinspection GoUnusedExportedFunction
func BindGormLogger(db gormDB) {
	db.SetLogger(log)
}

type gormDB interface {
	SetLogger(log gormLogger)
}

type gormLogger interface {
	Print(v ...interface{})
}
