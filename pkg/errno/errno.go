package errno

import "fmt"

/**
错误返回值格式：

{
"code": 10002,
"message": "Error occurred while binding the request body to the struct."
}
错误代码说明：
1: 系统级错误
00：服务模块代码
02：具体错误代码
服务级别错误：1 为系统级错误；2 为普通错误，通常是由用户非法操作引起的
服务模块为两位数：一个大型系统的服务模块通常不超过两位数，如果超过，说明这个系统该拆分了
错误码为两位数：防止一个模块定制过多的错误码，后期不好维护
code = 0 说明是正确返回，code > 0 说明是错误返回
错误通常包括系统级错误码和服务级错误码
建议代码中按服务模块将错误分类
错误码均为 >= 0 的数
在 apiserver 中 HTTP Code 固定为 http.StatusOK，错误码通过 code 来表示。
*/

type Errno struct {
	Code int
	Message string
}

func (err Errno) Error() string {
	return err.Message
}

type Err struct {
	Code int
	Message string
	Err error
}
func New(errno *Errno, err error) *Err {
	return &Err{
		Code:    errno.Code,
		Message: errno.Message,
		Err:     err,
	}
}

func (err *Err) Add(message string) error {
	err.Message += " " + message
	return err
}

func (err *Err) Addf(format string, args ...interface{}) error {
	err.Message += " " + fmt.Sprintf(format, args...)
	return err
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

func IsErrUserNotFound(err error) bool {
	code, _ := DecodeErr(err)
	return code == ErrUserNotFound.Code
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}
	return InternalServerError.Code, err.Error()
}