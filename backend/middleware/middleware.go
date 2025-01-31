package middleware

import "net/http"

const (
	HeaderKeyContentType  = "Content-Type"
	HeaderKeyContentValue = "application/json"
)

var Registry = []func(http.Handler) http.Handler{
	ContentTypeMiddleware,
}

func ContentTypeMiddleware(n http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set(HeaderKeyContentType, HeaderKeyContentValue)
		n.ServeHTTP(writer, request)
	})
}

func AddMiddleware(n http.Handler) http.Handler {
	for _, middleware := range Registry {
		n = middleware(n)
	}
	return n
}
