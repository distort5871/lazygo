package logger

import (
	"github.com/shiena/ansicolor"
	"io"
	"os"
	"sync"
	"time"
)

// brush is a color join function
type brush func(string) string

// newBrush return a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []brush{
	newBrush("1;37"), // Emergency          white
	newBrush("1;36"), // Alert              cyan
	newBrush("1;35"), // Critical           magenta
	newBrush("1;31"), // Error              red
	newBrush("1;33"), // Warning            yellow
	newBrush("1;32"), // Notice             green
	newBrush("1;34"), // Informational      blue
	newBrush("1;44"), // Debug              Background blue
}

type consoleLogWriter struct {
	sync.Mutex
	writer io.Writer
}

func newConsoleLogWriter(opt map[string]string) (logWriter, error) {
	cl := &consoleLogWriter{
		writer: ansicolor.NewAnsiColorWriter(os.Stdout),
	}
	return cl, nil
}

func (cl *consoleLogWriter) Writeln(b []byte, t time.Time) (int, error) {
	cl.Lock()
	defer cl.Unlock()
	hd, _, _ := formatTimeHeader(t)
	return cl.writer.Write(append(append(hd, b...), '\n'))
}

func (cl *consoleLogWriter) Close() error {
	return nil
}
func (cl *consoleLogWriter) Flush() error {
	return nil
}

func init() {
	// 注册适配器
	registry.add("console", adapterFunc(newConsoleLogWriter))
}