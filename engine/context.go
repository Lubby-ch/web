package engine

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	handlers   HandlersChain
	fullPath   string
	index      int8
	statusCode int
}

func (c *Context) ReadJson(req interface{}) error {
	request := c.Request
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, req)
}

func (c *Context) WriteJson(statusCode int, resp interface{}) error {
	c.Writer.WriteHeader(statusCode)
	respJson, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = c.Writer.Write(respJson)
	return err
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func NewContext(writer http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		Writer:  writer,
		Request: request,

		index: -1,
	}
}
