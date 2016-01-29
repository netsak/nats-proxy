package natsproxy

import "testing"

func TestPathVariable(t *testing.T) {
	req := &Request{
		URL: "/test/1234/tst",
	}
	resp := &Response{}
	ctx := newContext("/test/:token/:session", resp, req)

	tkn := ctx.PathVariable("token")
	if tkn != "1234" {
		t.Error("Defined path variable returned empty string")
	}

	unknwn := ctx.PathVariable("novalue")
	if unknwn != "" {
		t.Error("Non existing path variable returned non empty string")
	}

	unknwn = ctx.PathVariable("")
	if unknwn != "" {
		t.Error("Non existing path variable returned non empty string")
	}

	req.URL = ""
	tkn = ctx.PathVariable("token")
	if tkn != "" {
		t.Error("No variable in URL.Path returned non emtpy token")
	}

}