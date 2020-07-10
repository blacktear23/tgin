package tgin

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

var (
	dunno = []byte("???")
	cdot  = []byte("·")
	dot   = []byte(".")
	slash = []byte("/")
)

func RecoveryMiddleware(c *Context) {
	defer handlePanic(c)
	c.Next()
}

func handlePanic(c *Context) {
	if err := recover(); err != nil {
		brokenPipe := false
		if ne, ok := err.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}
		stack := dumpStack(3)
		log.Printf("[Recovery] panic recovered: %s\n%s\n", err, stack)
		if brokenPipe {
			c.Abort()
		} else {
			c.AbortWithStatus(500)
		}
	}
}

func dumpStack(skip int) []byte {
	buf := bytes.NewBuffer(nil)
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		fmt.Fprintf(buf, "\tfunc: %s\n", function(pc))
	}
	return buf.Bytes()
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, cdot, dot, -1)
	return name
}
