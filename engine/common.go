package engine

import "log"

var (
	default404Body = []byte("404 page not found")
	//default405Body = []byte("405 method not allowed")
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

func assert(condition bool, text string) {
	if !condition {
		panic(text)
	}
}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.Writer.WriteHeader(code)
	c.Next()
	_, err := c.Writer.Write(defaultMessage)
	if err != nil {
		log.Printf("cannot write message to writer during serve error: %v", err)
	}
}