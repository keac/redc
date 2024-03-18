package mod

import (
	"fmt"
	"os"
	"red-cloud/mod2"
	"time"
)

// RedcLog 将给定的消息记录到 "redc.log" 文件中。
func RedcLog(message string) {
	// 打开或创建日志文件
	file, err := os.OpenFile("redc.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	mod2.PrintOnError(err, "failed to open log file")
	defer file.Close()

	// 获取当前时间作为日志时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 构造日志消息并写入文件
	_, err = file.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
	mod2.PrintOnError(err, "failed to write to log file")

}
