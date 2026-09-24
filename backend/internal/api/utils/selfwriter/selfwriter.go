package selfwriter

import (
	"net/http"
)

type SelfWriter struct {
	http.ResponseWriter
	Code int
}

func (sw *SelfWriter) WriteHeader(code int) {
	sw.Code = code
	sw.ResponseWriter.WriteHeader(code)
}
