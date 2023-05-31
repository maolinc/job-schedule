package middleware

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestEr(t *testing.T) {
	err := errors.New("rpc error: code = Code(100001) desc = cron表达式错误")

	// 解析错误
	st, _ := status.FromError(err)
	code := st.Code()
	desc := st.Message()

	// 处理错误
	switch code {
	case codes.InvalidArgument:
		// 处理无效参数错误
		fmt.Println("Invalid Argument:", desc)
	case codes.NotFound:
		// 处理找不到资源错误
		fmt.Println("Not Found:", desc)
	default:
		// 处理其他错误
		fmt.Println("Unknown Error:", desc)
	}
}
