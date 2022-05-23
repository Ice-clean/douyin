package middleware

import (
	"github.com/RaymondCode/simple-demo/result"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}

func Recovery() gin.HandlerFunc {
	return RecoveryWithWriter(gin.DefaultErrorWriter)
}

func RecoveryWithWriter(out io.Writer) gin.HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 有异常
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				if logger != nil {
					// 打印真正的错误栈
					debug.PrintStack()
					logger.Printf("程序发生了panic: %v\n", err)
					logger.Printf("恢复了 %s panic recovered:\n%s", timeFormat(time.Now()), err)
				}
				Gin := result.Gin{
					C: c,
				}
				// 如果连接已断开，我们无法向其写入状态。
				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
				} else {
					// 向前端返回错误信息，但不展示真正的错误
					Gin.AbortWithStatusJSON()
				}
			}
		}()
		// 函数继续调用
		c.Next()
	}
}
