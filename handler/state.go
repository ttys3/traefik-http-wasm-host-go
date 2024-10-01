package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/http-wasm/http-wasm-host-go/api/handler"
)

// requestStateKey is a context.Context value associated with a requestState
// pointer to the current request.
type requestStateKey struct{}

func requestStateFromContext(ctx context.Context) *requestState {
	return ctx.Value(requestStateKey{}).(*requestState)
}

type requestState struct {
	afterNext          bool
	requestBodyReader  io.ReadCloser
	requestBodyWriter  io.Writer
	responseBodyReader io.ReadCloser
	responseBodyWriter io.Writer

	// features are the current request's features which may be more than
	// Middleware.Features.
	features handler.Features

	putPool func(x any)
	g       *guest
}

func (r *requestState) closeRequest() (err error) {
	if reqBW := r.requestBodyWriter; reqBW != nil {
		if f, ok := reqBW.(http.Flusher); ok {
			f.Flush()
		}
		r.requestBodyWriter = nil
	}
	if reqBR := r.requestBodyReader; reqBR != nil {
		err = reqBR.Close()
		r.requestBodyReader = nil
	}
	return
}

// Close releases all resources for the current request, including:
//   - putting the guest module back into the pool
//   - releasing any request body resources
//   - releasing any response body resources
func (r *requestState) Close() (err error) {
	if g := r.g; g != nil {
		fmt.Printf("=================== xxoo putPool pool put back: requestState.Close\n")
		r.putPool(r.g)
		r.g = nil
	} else {
		fmt.Printf("=================== xxoo fuck r.g is nil: requestState.Close\n")
	}
	err = r.closeRequest()
	if respBW := r.responseBodyWriter; respBW != nil {
		if f, ok := respBW.(http.Flusher); ok {
			f.Flush()
		}
		r.responseBodyWriter = nil
	}
	if respBR := r.responseBodyReader; respBR != nil {
		err = respBR.Close()
		r.responseBodyReader = nil
	}
	return
}
