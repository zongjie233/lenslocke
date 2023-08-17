package errors

// 包装错误
func Public(err error, msg string) error {
	return nil
}

type publicError struct {
	err error
	msg string
}

func (pe publicError) Error() string {
	return pe.err.Error()
}
func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}
