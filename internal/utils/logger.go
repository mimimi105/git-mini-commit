package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

// LogLevel ログレベル
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String ログレベルの文字列表現
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger 構造化ログ機能
type Logger struct {
	level    LogLevel
	output   io.Writer
	prefix   string
	context  map[string]interface{}
}

// NewLogger 新しいLoggerインスタンスを作成
func NewLogger(level LogLevel, output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}
	
	return &Logger{
		level:   level,
		output:  output,
		prefix:  "",
		context: make(map[string]interface{}),
	}
}

// WithPrefix プレフィックスを設定したLoggerを作成
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		level:   l.level,
		output:  l.output,
		prefix:  prefix,
		context: l.context,
	}
}

// WithContext コンテキストを設定したLoggerを作成
func (l *Logger) WithContext(key string, value interface{}) *Logger {
	newContext := make(map[string]interface{})
	for k, v := range l.context {
		newContext[k] = v
	}
	newContext[key] = value
	
	return &Logger{
		level:   l.level,
		output:  l.output,
		prefix:  l.prefix,
		context: newContext,
	}
}

// SetLevel ログレベルを設定
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetOutput 出力先を設定
func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

// log 内部ログ関数
func (l *Logger) log(level LogLevel, message string, args ...interface{}) {
	if level < l.level {
		return
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMessage := fmt.Sprintf(message, args...)
	
	var logLine string
	if l.prefix != "" {
		logLine = fmt.Sprintf("[%s] %s %s: %s", timestamp, level.String(), l.prefix, formattedMessage)
	} else {
		logLine = fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), formattedMessage)
	}
	
	// コンテキスト情報を追加
	if len(l.context) > 0 {
		logLine += " | Context:"
		for key, value := range l.context {
			logLine += fmt.Sprintf(" %s=%v", key, value)
		}
	}
	
	fmt.Fprintln(l.output, logLine)
}

// Debug デバッグログを出力
func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(DEBUG, message, args...)
}

// Info 情報ログを出力
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(INFO, message, args...)
}

// Warn 警告ログを出力
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(WARN, message, args...)
}

// Error エラーログを出力
func (l *Logger) Error(message string, args ...interface{}) {
	l.log(ERROR, message, args...)
}

// Debugf デバッグログをフォーマット付きで出力
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(format, args...)
}

// Infof 情報ログをフォーマット付きで出力
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(format, args...)
}

// Warnf 警告ログをフォーマット付きで出力
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(format, args...)
}

// Errorf エラーログをフォーマット付きで出力
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(format, args...)
}

// LogError エラーをログに記録
func (l *Logger) LogError(err error, message string, args ...interface{}) {
	if err != nil {
		l.Error("%s: %v", fmt.Sprintf(message, args...), err)
	}
}

// LogErrorf エラーをフォーマット付きでログに記録
func (l *Logger) LogErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		l.Errorf("%s: %v", fmt.Sprintf(format, args...), err)
	}
}

// GetLevel 現在のログレベルを取得
func (l *Logger) GetLevel() LogLevel {
	return l.level
}

// IsDebugEnabled デバッグログが有効かチェック
func (l *Logger) IsDebugEnabled() bool {
	return l.level <= DEBUG
}

// IsInfoEnabled 情報ログが有効かチェック
func (l *Logger) IsInfoEnabled() bool {
	return l.level <= INFO
}

// IsWarnEnabled 警告ログが有効かチェック
func (l *Logger) IsWarnEnabled() bool {
	return l.level <= WARN
}

// IsErrorEnabled エラーログが有効かチェック
func (l *Logger) IsErrorEnabled() bool {
	return l.level <= ERROR
}

