package errorx

import "fmt"

const (
	RPC_ERROR = 1
)

type CodeError struct {
	errCode uint32
	errMsg  string
}

// 返回给前端的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d, ErrMsg:%s", e.errCode, e.errMsg)
}

func NewMsg(errMsg string) *CodeError {
	return &CodeError{errCode: RPC_ERROR, errMsg: errMsg}
}
