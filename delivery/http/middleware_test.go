package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/zitryss/perfmon/internal/testing"
)

func TestNoDirListing(t *testing.T) {
	type args struct {
		method string
		target string
		body   io.Reader
	}
	type want struct {
		code    int
		content string
		body    string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{
				method: "GET",
				target: "/static/imgs/1.png",
				body:   nil,
			},
			want: want{
				code:    200,
				content: "",
				body:    "",
			},
		},
		{
			name: "",
			args: args{
				method: "GET",
				target: "/static/imgs/",
				body:   nil,
			},
			want: want{
				code:    404,
				content: "text/plain; charset=utf-8",
				body:    "404 page not found\n",
			},
		},
	}
	fn := func(w http.ResponseWriter, r *http.Request) {}
	h := noDirListing(http.HandlerFunc(fn))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.target, tt.args.body)
			h.ServeHTTP(w, r)
			CheckStatusCode(t, w, tt.want.code)
			CheckContentType(t, w, tt.want.content)
			CheckBody(t, w, tt.want.body)
		})
	}
}

func TestRecovery(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) { panic("") }
	h := recovery(http.HandlerFunc(fn))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, r)
	CheckStatusCode(t, w, 500)
	CheckContentType(t, w, "text/plain; charset=utf-8")
	CheckBody(t, w, "Internal Server Error\n")
}
