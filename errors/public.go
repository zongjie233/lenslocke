package errors

// 包装错误
func Public(err error, msg string) error {
	// 返回一个错误
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

// Error 方法返回错误的字符串
func (pe publicError) Error() string {
	// 返回错误的字符串
	return pe.err.Error()
}

// Public 方法返回错误的描述
func (pe publicError) Public() string {
	// 返回错误的描述
	return pe.msg
}

// Unwrap 方法返回错误的原始错误
func (pe publicError) Unwrap() error {
	// 返回错误的原始错误
	return pe.err
}
