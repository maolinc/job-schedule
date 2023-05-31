package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
)

type LogMiddleware struct {
	// 可以在这里添加操作数据的依赖
	logx.Logger
}

func NewLogMiddleware() *LogMiddleware {
	return &LogMiddleware{
		logx.WithContext(context.Background()),
	}
}

type logResponseWriter struct {
	writer http.ResponseWriter
	code   int
	buf    *bytes.Buffer
}

func newLogResponseWriter(writer http.ResponseWriter, code int) *logResponseWriter {
	var buf bytes.Buffer
	return &logResponseWriter{
		writer: writer,
		code:   code,
		buf:    &buf,
	}
}

func (w *logResponseWriter) Write(bs []byte) (int, error) {
	w.buf.Write(bs)
	return w.writer.Write(bs)
}

func (w *logResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *logResponseWriter) WriteHeader(code int) {
	w.code = code
	w.writer.WriteHeader(code)
}

func (m *LogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dup io.ReadCloser
		dup, r.Body, _ = drainBody(r.Body)
		lwr := newLogResponseWriter(w, http.StatusOK)
		next(lwr, r)
		r.Body = dup
		m.logDetailLogic(r, lwr)
	}
}

type status struct {
	ResultStatus struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}
}

func (m *LogMiddleware) logDetailLogic(request *http.Request, response *logResponseWriter) {
	s := status{}
	err := json.Unmarshal(response.buf.Bytes(), &s)
	if err == nil {
		if s.ResultStatus.Code != "0" {
			bs, _ := io.ReadAll(request.Body)
			m.Errorf("req: %s, resp: %+v", string(bs), s.ResultStatus)
		}
	}
}

// drainBody from httputil.drainBody
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
