package utils

import (
	"fmt"
	"strings"
)

// ErrorType エラータイプ
type ErrorType int

const (
	// FileOperationError ファイル操作エラー
	FileOperationError ErrorType = iota
	// GitOperationError Git操作エラー
	GitOperationError
	// StorageError ストレージエラー
	StorageError
	// ValidationError バリデーションエラー
	ValidationError
	// ConfigurationError 設定エラー
	ConfigurationError
	// NetworkError ネットワークエラー
	NetworkError
	// PermissionError 権限エラー
	PermissionError
	// UnknownError 不明なエラー
	UnknownError
)

// String エラータイプの文字列表現
func (et ErrorType) String() string {
	switch et {
	case FileOperationError:
		return "FileOperation"
	case GitOperationError:
		return "GitOperation"
	case StorageError:
		return "Storage"
	case ValidationError:
		return "Validation"
	case ConfigurationError:
		return "Configuration"
	case NetworkError:
		return "Network"
	case PermissionError:
		return "Permission"
	case UnknownError:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// AppError アプリケーション固有のエラー
type AppError struct {
	Type      ErrorType
	Message   string
	Context   map[string]interface{}
	Cause     error
	Timestamp string
}

// Error エラーメッセージを返す
func (e *AppError) Error() string {
	var parts []string
	
	// エラータイプ
	parts = append(parts, fmt.Sprintf("[%s]", e.Type.String()))
	
	// メッセージ
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	
	// コンテキスト情報
	if len(e.Context) > 0 {
		var contextParts []string
		for key, value := range e.Context {
			contextParts = append(contextParts, fmt.Sprintf("%s=%v", key, value))
		}
		parts = append(parts, fmt.Sprintf("Context: %s", strings.Join(contextParts, ", ")))
	}
	
	// 原因エラー
	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("Cause: %v", e.Cause))
	}
	
	return strings.Join(parts, " | ")
}

// Unwrap 原因エラーを返す
func (e *AppError) Unwrap() error {
	return e.Cause
}

// ErrorHandler エラーハンドリング機能
type ErrorHandler struct {
	logger *Logger
}

// NewErrorHandler 新しいErrorHandlerインスタンスを作成
func NewErrorHandler(logger *Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// NewError 新しいAppErrorを作成
func (eh *ErrorHandler) NewError(errorType ErrorType, message string, cause error) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
		Context: make(map[string]interface{}),
	}
}

// NewErrorWithContext コンテキスト付きのAppErrorを作成
func (eh *ErrorHandler) NewErrorWithContext(errorType ErrorType, message string, cause error, context map[string]interface{}) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
		Context: context,
	}
}

// WrapError 既存のエラーをAppErrorでラップ
func (eh *ErrorHandler) WrapError(errorType ErrorType, message string, cause error) *AppError {
	return eh.NewError(errorType, message, cause)
}

// WrapErrorWithContext コンテキスト付きで既存のエラーをAppErrorでラップ
func (eh *ErrorHandler) WrapErrorWithContext(errorType ErrorType, message string, cause error, context map[string]interface{}) *AppError {
	return eh.NewErrorWithContext(errorType, message, cause, context)
}

// HandleError エラーを処理してログに記録
func (eh *ErrorHandler) HandleError(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	
	formattedMessage := fmt.Sprintf(message, args...)
	
	// AppErrorの場合はそのまま返す
	if appErr, ok := err.(*AppError); ok {
		eh.logger.LogError(appErr, formattedMessage)
		return appErr
	}
	
	// 通常のエラーの場合はラップ
	wrappedErr := eh.WrapError(UnknownError, formattedMessage, err)
	eh.logger.LogError(wrappedErr, formattedMessage)
	return wrappedErr
}

// HandleFileError ファイル操作エラーを処理
func (eh *ErrorHandler) HandleFileError(err error, operation, filepath string) error {
	if err == nil {
		return nil
	}
	
	context := map[string]interface{}{
		"operation": operation,
		"filepath":  filepath,
	}
	
	return eh.WrapErrorWithContext(FileOperationError, fmt.Sprintf("file operation failed: %s", operation), err, context)
}

// HandleGitError Git操作エラーを処理
func (eh *ErrorHandler) HandleGitError(err error, operation string) error {
	if err == nil {
		return nil
	}
	
	context := map[string]interface{}{
		"operation": operation,
	}
	
	return eh.WrapErrorWithContext(GitOperationError, fmt.Sprintf("git operation failed: %s", operation), err, context)
}

// HandleStorageError ストレージエラーを処理
func (eh *ErrorHandler) HandleStorageError(err error, operation string) error {
	if err == nil {
		return nil
	}
	
	context := map[string]interface{}{
		"operation": operation,
	}
	
	return eh.WrapErrorWithContext(StorageError, fmt.Sprintf("storage operation failed: %s", operation), err, context)
}

// HandleValidationError バリデーションエラーを処理
func (eh *ErrorHandler) HandleValidationError(err error, field, value string) error {
	if err == nil {
		return nil
	}
	
	context := map[string]interface{}{
		"field": field,
		"value": value,
	}
	
	return eh.WrapErrorWithContext(ValidationError, fmt.Sprintf("validation failed for field: %s", field), err, context)
}

// HandlePermissionError 権限エラーを処理
func (eh *ErrorHandler) HandlePermissionError(err error, resource string) error {
	if err == nil {
		return nil
	}
	
	context := map[string]interface{}{
		"resource": resource,
	}
	
	return eh.WrapErrorWithContext(PermissionError, fmt.Sprintf("permission denied for resource: %s", resource), err, context)
}

// LogError エラーをログに記録
func (eh *ErrorHandler) LogError(err error, message string, args ...interface{}) {
	if err != nil {
		eh.logger.LogError(err, message, args...)
	}
}

// LogErrorf エラーをフォーマット付きでログに記録
func (eh *ErrorHandler) LogErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		eh.logger.LogErrorf(err, format, args...)
	}
}

// IsErrorType エラーが特定のタイプかチェック
func (eh *ErrorHandler) IsErrorType(err error, errorType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errorType
	}
	return false
}

// GetErrorType エラーのタイプを取得
func (eh *ErrorHandler) GetErrorType(err error) ErrorType {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type
	}
	return UnknownError
}

