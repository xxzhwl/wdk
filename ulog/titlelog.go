package ulog

// TitleLog 带标题的Log
type TitleLog struct {
	Title string
}

// NewTitleLog 生成一个带标题的log
func NewTitleLog(title string) TitleLog {
	return TitleLog{title}
}

// Info 信息日志
func (t TitleLog) Info(msg string) {
	Info(t.Title, msg)
}

// InfoF 带模板的信息日志
func (t TitleLog) InfoF(template string, args ...any) {
	InfoF(t.Title, template, args)
}

// Warn 警告日志
func (t TitleLog) Warn(msg string) {
	Info(t.Title, msg)
}

// WarnF 带模板的警告日志
func (t TitleLog) WarnF(template string, args ...any) {
	WarnF(t.Title, template, args)
}

// Error 错误日志
func (t TitleLog) Error(msg string) {
	Info(t.Title, msg)
}

// ErrorF 带模板的错误日志
func (t TitleLog) ErrorF(template string, args ...any) {
	ErrorF(t.Title, template, args)
}
