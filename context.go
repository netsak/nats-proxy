package natsproxy

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var removeQueryRxp = regexp.MustCompile("[?]{1}.*")

// Context wraps the
// processed request/response
type Context struct {
	Request    *Request
	Response   *Response
	index      int
	abortIndex int
	params     map[string]int
}

// IsAborted returns true
// if the request in context
// were aborted by previous
// middleware
func (c *Context) IsAborted() bool {
	return c.index >= c.abortIndex
}

// Abort abortsthe
// request that it won's be
// processed further
func (c *Context) Abort() {
	c.abortIndex = c.index
}

// AbortWithJSON aborts the request
// and sets the HTTP status code to 500.
func (c *Context) AbortWithJSON(obj interface{}) {
	c.Abort()
	c.JSON(500, obj)
}

// BindJSON unmarshall the
// request body to given
// struct
func (c *Context) BindJSON(obj interface{}) error {
	if err := json.Unmarshal(c.Request.Body, obj); err != nil {
		return err
	}
	return nil
}

// JSON writes the serialized
// json to response
func (c *Context) JSON(statusCode int, obj interface{}) {
	c.Response.StatusCode = statusCode
	bytes, err := json.Marshal(obj)
	if err != nil {
		c.writeError(err)
	}
	c.Response.Body = bytes

}

// PathVariable returns
// the path variable
// based on its name (:xxx) defined
// in subscription URL
func (c *Context) PathVariable(name string) string {
	URL := removeQueryRxp.ReplaceAllString(c.Request.URL, "")
	pathParams := strings.Split(URL, "/")

	index, ok := c.params[name]
	if !ok {
		return ""
	}
	if len(pathParams) <= index {
		return ""
	}
	return pathParams[index]
}

// FormVariable returns the
// variable from request form if
// available.
func (c *Context) FormVariable(name string) string {
	return c.Request.Form.Get(name)
}

func (c *Context) writeError(err error) {
	c.Response.StatusCode = 500
	c.Response.Body = []byte(err.Error())
}

func newContext(url string, res *Response, req *Request) *Context {
	m := buildParamMap(url)
	return &Context{
		req,
		res,
		0,
		1<<31 - 1,
		m,
	}
}

func buildParamMap(url string) map[string]int {
	m := make(map[string]int)
	prmArr := strings.Split(url, "/")
	for i, prm := range prmArr {
		if len(prm) > 0 && prm[:1] == ":" {
			m[prm[1:]] = i
		}
	}
	return m
}

func getPathVariableAtPlace(url string, place int) (string, error) {
	parsedPath := strings.Split(url[1:], "/")
	if len(parsedPath) < place {
		return "", fmt.Errorf("Variable not found")
	}
	return parsedPath[place], nil
}
