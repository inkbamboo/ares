// Package logger 提供自定义彩色日志器实现（独立于 logrus）。
package logger

// Import packages
import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
	"time"
)

var (
	// colors 存储不同日志级别对应的颜色代码
	colors map[string]string

	// logNo 日志序号计数器
	logNo uint64
)

// Color numbers for stdout
const (
	Black = (iota + 30)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Worker 是底层日志工作者，负责实际的日志输出和颜色处理。
type Worker struct {
	Minion *log.Logger
	Color  int
}

// Info 包含单条日志的所有信息。
type Info struct {
	Id      uint64
	Time    string
	Module  string
	Level   string
	Message string
	format  string
}

// Logger 是面向用户的日志接口，提供模块化的日志记录功能。
type Logger struct {
	Module string
	worker *Worker
}

// Output 返回格式化后的日志字符串。
func (r *Info) Output() string {

	var idString string
	if r.Id < 10 {
		idString = fmt.Sprintf("00%d", r.Id)
	} else if r.Id >= 10 && r.Id <= 99 {
		idString = fmt.Sprintf("0%d", r.Id)
	} else {
		idString = fmt.Sprintf("%d", r.Id)
	}

	msg := fmt.Sprintf(r.format, idString, r.Time, r.Level, r.Message)
	return msg
}

// NewWorker 创建一个新的 Worker 实例。
func NewWorker(prefix string, flag int, color int, out io.Writer) *Worker {
	return &Worker{Minion: log.New(out, prefix, flag), Color: color}
}

// Log 根据日志级别输出带颜色或不带颜色的日志。
func (w *Worker) Log(level string, calldepth int, info *Info) error {
	if w.Color != 0 {
		buf := &bytes.Buffer{}
		buf.Write([]byte(colors[level]))
		buf.Write([]byte(info.Output()))
		buf.Write([]byte("\033[0m"))
		return w.Minion.Output(calldepth+1, buf.String())
	} else {
		return w.Minion.Output(calldepth+1, info.Output())
	}
}

// colorString 返回指定颜色的 ANSI 转义序列字符串。
func colorString(color int) string {
	return fmt.Sprintf("\033[%dm", int(color))
}

// initColors 初始化日志级别颜色映射表。
func initColors() {
	colors = map[string]string{
		"CRITICAL": colorString(Magenta),
		"ERROR":    colorString(Red),
		"WARNING":  colorString(Yellow),
		"NOTICE":   colorString(Green),
		"DEBUG":    colorString(Cyan),
		"INFO":     colorString(White),
	}
}

// NewLogger 创建一个新的 Logger 实例。
// 可选参数：module(string)、color(int)、out(io.Writer)。
func NewLogger(args ...interface{}) (*Logger, error) {
	initColors()

	var module = "DEFAULT"
	var color = 1
	var out io.Writer = os.Stderr

	for _, arg := range args {
		switch t := arg.(type) {
		case string:
			module = t
		case int:
			color = t
		case io.Writer:
			out = t
		default:
			panic("logger: Unknown argument")
		}
	}
	newWorker := NewWorker("", 0, color, out)
	return &Logger{Module: module, worker: newWorker}, nil
}

// Log 记录指定级别和消息的日志。
func (l *Logger) Log(lvl string, message string) {
	var formatString string = "#%s %s ▶ %.3s %s"
	info := &Info{
		Id:      atomic.AddUint64(&logNo, 1),
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Module:  l.Module,
		Level:   lvl,
		Message: message,
		format:  formatString,
	}
	l.worker.Log(lvl, 2, info)
}

// Fatal 记录 CRITICAL 级别日志后退出程序。
func (l *Logger) Fatal(message string) {
	l.Log("CRITICAL", message)
	os.Exit(1)
}

// Panic 记录 CRITICAL 级别日志后触发 panic。
func (l *Logger) Panic(message string) {
	l.Log("CRITICAL", message)
	panic(message)
}

// Critical 记录 CRITICAL 级别日志。
func (l *Logger) Critical(message string) {
	l.Log("CRITICAL", message)
}

// Error 记录 ERROR 级别日志。
func (l *Logger) Error(message string) {
	l.Log("ERROR", message)
}

// Warning 记录 WARNING 级别日志。
func (l *Logger) Warning(message string) {
	l.Log("WARNING", message)
}

// Notice 记录 NOTICE 级别日志。
func (l *Logger) Notice(message string) {
	l.Log("NOTICE", message)
}

// Info 记录 INFO 级别日志。
func (l *Logger) Info(message string) {
	l.Log("INFO", message)
}

// Debug 记录 DEBUG 级别日志。
func (l *Logger) Debug(message string) {
	l.Log("DEBUG", message)
}
