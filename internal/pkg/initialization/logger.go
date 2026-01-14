package initialization

import (
	"fmt"
	"os"
	"path/filepath"
	"shunshun/internal/pkg/global"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// lazyWriteSyncer 延迟打开文件的WriteSyncer
type lazyWriteSyncer struct {
	filePath string
	file     *os.File
	mutex    sync.Mutex
}

// NewLazyWriteSyncer 创建一个延迟打开文件的WriteSyncer
func NewLazyWriteSyncer(filePath string) *lazyWriteSyncer {
	return &lazyWriteSyncer{
		filePath: filePath,
	}
}

// Write 实现WriteSyncer接口
func (l *lazyWriteSyncer) Write(p []byte) (n int, err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file == nil {
		// 第一次写入时创建文件
		file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		l.file = file
	}

	return l.file.Write(p)
}

// Sync 实现WriteSyncer接口
func (l *lazyWriteSyncer) Sync() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.file != nil {
		return l.file.Sync()
	}
	return nil
}

// InitLogger 初始化zap日志
func InitLogger() *zap.Logger {
	// 获取当前文件的目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前工作目录失败: %v\n", err)
		os.Exit(1)
	}
	// 项目根目录（相对于当前工作目录）
	projectRoot := filepath.Join(currentDir, "..", "..")

	// 从配置中获取日志目录，默认为internal/logger
	if global.AppConf == nil {
		fmt.Printf("global.AppConf 为 nil\n")
		os.Exit(1)
	}
	logDir := global.AppConf.Zap.LogDir
	if logDir == "" {
		logDir = filepath.Join(projectRoot, "internal", "logger")
	} else {
		// 如果配置的是相对路径，则相对于项目根目录
		if !filepath.IsAbs(logDir) {
			logDir = filepath.Join(projectRoot, logDir)
		}
	}

	// 转换为绝对路径
	logDir, err = filepath.Abs(logDir)
	if err != nil {
		fmt.Printf("转换日志目录为绝对路径失败: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
		os.Exit(1)
	}

	// 从配置中获取保留天数，默认为7天
	maxAge := global.AppConf.Zap.MaxAge
	if maxAge <= 0 {
		maxAge = 7
	}

	// 清除过期日志
	cleanupOldLogs(logDir, maxAge)
	// 获取当前日期
	today := time.Now().Format("2006-01-02")
	todayLogDir := filepath.Join(logDir, today)
	if err := os.MkdirAll(todayLogDir, 0755); err != nil {
		fmt.Printf("创建今日日志目录失败: %v\n", err)
		os.Exit(1)
	}
	// 定义日志文件路径
	infoLogPath := filepath.Join(todayLogDir, "info.log")
	errorLogPath := filepath.Join(todayLogDir, "error.log")
	// 创建 info 日志文件
	infoFile, err := os.OpenFile(infoLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("创建info日志文件失败: %v\n", err)
		os.Exit(1)
	}
	// 配置 zap
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	// 创建编码器
	infoEncoder := zapcore.NewConsoleEncoder(config)
	errorEncoder := zapcore.NewConsoleEncoder(config)
	// 创建核心
	// Info核心只处理Info级别
	infoCore := zapcore.NewCore(infoEncoder, zapcore.AddSync(infoFile), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	}))
	// Error核心处理Error及以上级别
	// 使用延迟写入器，只有在有错误时才创建error.log文件
	errorCore := zapcore.NewCore(errorEncoder, NewLazyWriteSyncer(errorLogPath), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	}))
	// 创建复合核心
	core := zapcore.NewTee(
		infoCore,
		errorCore,
	)
	// logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger
}

// cleanupOldLogs 清除过期日志
func cleanupOldLogs(logDir string, maxAge int) {
	// 获取过期日期
	expireDate := time.Now().AddDate(0, 0, -maxAge)
	expireDateStr := expireDate.Format("2006-01-02")

	// 读取日志目录
	dirs, err := os.ReadDir(logDir)
	if err != nil {
		fmt.Printf("读取日志目录失败: %v\n", err)
		return
	}

	// 遍历目录，删除7天前的日志
	for _, dir := range dirs {
		if dir.IsDir() {
			// 检查目录名是否为日期格式
			if _, err := time.Parse("2006-01-02", dir.Name()); err == nil {
				// 如果目录日期早于7天前，则删除
				if dir.Name() < expireDateStr {
					dirPath := filepath.Join(logDir, dir.Name())
					if err := os.RemoveAll(dirPath); err != nil {
						fmt.Printf("删除旧日志目录失败: %v\n", err)
					}
				}
			}
		}
	}
}
