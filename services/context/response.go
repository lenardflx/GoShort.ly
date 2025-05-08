package context

import (
	"net/http"
)

// ResponseWriter defines the interface used by our context-aware writer
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher

	Before(fn func(ResponseWriter))
	WrittenSize() int
	WrittenStatus() int
}

// Response wraps http.ResponseWriter to track status, written bytes, and allow hooks
type Response struct {
	http.ResponseWriter
	written        int
	status         int
	beforeFuncs    []func(ResponseWriter)
	beforeExecuted bool
}

// WrapResponseWriter wraps a standard http.ResponseWriter into our Response
func WrapResponseWriter(w http.ResponseWriter) *Response {
	if resp, ok := w.(*Response); ok {
		return resp
	}
	return &Response{ResponseWriter: w}
}

// Write implements http.ResponseWriter.Write
func (r *Response) Write(b []byte) (int, error) {
	r.executeBefore()
	size, err := r.ResponseWriter.Write(b)
	r.written += size
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return size, err
}

// WriteHeader implements http.ResponseWriter.WriteHeader
func (r *Response) WriteHeader(code int) {
	r.executeBefore()
	if r.status == 0 {
		r.status = code
		r.ResponseWriter.WriteHeader(code)
	}
}

// Flush implements http.Flusher
func (r *Response) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Before registers a function to be called before the response is written
func (r *Response) Before(fn func(ResponseWriter)) {
	r.beforeFuncs = append(r.beforeFuncs, fn)
}

// WrittenSize returns the number of bytes written
func (r *Response) WrittenSize() int {
	return r.written
}

// WrittenStatus returns the HTTP status code written
func (r *Response) WrittenStatus() int {
	return r.status
}

func (r *Response) executeBefore() {
	if !r.beforeExecuted {
		for _, fn := range r.beforeFuncs {
			fn(r)
		}
		r.beforeExecuted = true
	}
}

func (r *Response) Written() bool {
	return r.written > 0
}
