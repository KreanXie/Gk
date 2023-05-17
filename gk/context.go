package gk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H is an alias of map[string]interface{} which stores and passes Header of http request
type H map[string]interface{}

// Context stores the messages which comes from middleware and is strongly related with the current request
// Also Context can store parameter of dynamic route
// In one word, Context is a message set
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int

	// middleware
	handlers []HandlerFunc
	index    int

	// engine pointer
	engine *Engine
}

// New returns a pointer of new Context object
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	n := len(c.handlers)
	for ; c.index < n; c.index++ {
		c.handlers[c.index](c)
	}
}

// Fail is a common use function for err handling in gk
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm is the encapsulation of FormValue() of http.Request
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query is the encapsulation of Query() of http.Request.URL
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status sets StatusCode with receving code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader is the encapsulation of Header() of http.ResponseWriter
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String can make constructing String much faster
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON makes constructing json faster
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data can make constructing Data much faster
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML can make constructing HTML much faster
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data)
	if err != nil {
		c.Fail(500, err.Error())
	}
}
