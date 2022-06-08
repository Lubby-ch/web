package engine

import "net/http"

type MapBasedHandler struct {
	Handlers map[string]func(*Context)
}

func (m *MapBasedHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := m.key(request.Method, request.URL.Path)
	if handler, ok := m.Handlers[key]; ok {
		handler(NewContext(writer, request))
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("page not found"))
	}
}

func (m *MapBasedHandler) key(method, pattern string) string {
	return method + pattern
}
