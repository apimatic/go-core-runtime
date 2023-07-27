package https

import (
	"net/http"
)

// HttpCallExecutor is a function type that represents the execution of an HTTP call and returns the HTTP context.
type HttpCallExecutor func(request *http.Request) HttpContext

// HttpInterceptor is a function type that represents an HTTP interceptor, which intercepts and processes an HTTP call.
type HttpInterceptor func(request *http.Request, next HttpCallExecutor) HttpContext

// PassThroughInterceptor is an HTTP interceptor that passes the request to the next HttpCallExecutor in the chain.
// It does not modify the request or response and acts as a no-operation interceptor.
func PassThroughInterceptor(
	req *http.Request,
	next HttpCallExecutor,
) HttpContext {
	return next(req)
}

// CallHttpInterceptors chains multiple HTTP interceptors together to create a composite HttpCallExecutor.
// The interceptors are applied in the reverse order of their appearance in the interceptors slice.
// Each interceptor gets a chance to process the request before passing it to the next interceptor in the chain,
// and finally, the original HttpCallExecutor is called to execute the HTTP call.
// The composite HttpCallExecutor is returned, which includes the logic of all the interceptors.
func CallHttpInterceptors(
	interceptors []HttpInterceptor,
	client HttpCallExecutor,
) HttpCallExecutor {
	var next HttpCallExecutor = client
	for index := len(interceptors) - 1; index >= 0; index-- {
		current := interceptors[index]
		last := next
		next = func(req *http.Request) HttpContext {
			return current(req, last)
		}
	}
	return next
}
