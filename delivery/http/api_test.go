package http

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/zitryss/perfmon/domain"
	"github.com/zitryss/perfmon/infrastructure/database"
	. "github.com/zitryss/perfmon/internal/testing"
)

var (
	dbMem = database.NewMem()
)

func TestApiUploadMem(t *testing.T) {
	type args struct {
		use    Usecaser
		method string
		target string
		body   string
	}
	type want struct {
		code    int
		content string
		body    string
	}
	useOk := NewUsecase(dbMem)
	useFail := usecaseFail{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p1","version":"v1","attributes":["a11","a12"],"name":"n1","measurement":"m1","timestamp":"2018-11-08T16:23:50Z","value":1}`,
			},
			want: want{
				code:    201,
				content: "",
				body:    "",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","timestamp":"2017-10-29T16:18:07Z","value":48}`,
			},
			want: want{
				code:    201,
				content: "",
				body:    "",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62","a63"],"name":"n6","measurement":"m6","timestamp":"2017-10-29T16:18:07Z","value":48}`,
			},
			want: want{
				code:    409,
				content: "text/plain; charset=utf-8",
				body:    "Resource Already Exists\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p1`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"version":"v1","attributes":["a11","a12"],"name":"n1","measurement":"m1","timestamp":"2006-01-02T15:04:05Z","value":1}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p1","attributes":["a11","a12"],"name":"n1","measurement":"m1","timestamp":"2006-01-02T15:04:05Z","value":1}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p1","version":"v1","attributes":["a11","a12"],"measurement":"m1","timestamp":"2006-01-02T15:04:05Z","value":1}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p1","version":"v1","attributes":["a11","a12"],"name":"n1","timestamp":"2006-01-02T15:04:05Z","value":1}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p1","version":"v1","attributes":["a11","a12"],"name":"n1","measurement":"m1","value":1}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "POST",
				target: "/",
				body:   `{"product":"p1","version":"v1","attributes":["a11","a12"],"name":"n1","measurement":"m1","timestamp":"2006-01-02T15:04:05Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useFail,
				method: "POST",
				target: "/",
				body:   `{"product":"p1","version":"v1","attributes":["a11","a12"],"name":"n1","measurement":"m1","timestamp":"2006-01-02T15:04:05Z","value":1}`,
			},
			want: want{
				code:    500,
				content: "text/plain; charset=utf-8",
				body:    "Internal Server Error\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(tt.args.use)
			fn := api.upload()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.target, strings.NewReader(tt.args.body))
			fn(w, r, nil)
			CheckStatusCode(t, w, tt.want.code)
			CheckContentType(t, w, tt.want.content)
			CheckBody(t, w, tt.want.body)
		})
	}
	job := &Job{
		Product:     "p1",
		Version:     "v1",
		Attributes:  []string{"a11", "a12"},
		Name:        "n1",
		Measurement: "m1",
		Timestamp:   time.Date(2018, 11, 8, 16, 23, 50, 0, time.UTC),
		Value:       1,
	}
	err := dbMem.Delete(job)
	if err != nil {
		t.Errorf("error occurred: %+v", err)
	}
	job = &Job{
		Product:     "p6",
		Version:     "v6",
		Attributes:  []string{"a61", "a62"},
		Name:        "n6",
		Measurement: "m6",
		Timestamp:   time.Date(2017, 10, 29, 16, 18, 07, 0, time.UTC),
		Value:       48,
	}
	err = dbMem.Delete(job)
	if err != nil {
		t.Errorf("error occurred: %+v", err)
	}
}

func TestApiProductsMem(t *testing.T) {
	type args struct {
		use    Usecaser
		method string
		target string
		body   string
	}
	type want struct {
		code    int
		content string
		body    string
	}
	useOk := NewUsecase(dbMem)
	useFail := usecaseFail{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/products/",
				body:   "",
			},
			want: want{
				code:    200,
				content: "application/json; charset=utf-8",
				body:    `{"products":["p2","p3","p4","p5","p6","p7","pY"]}` + "\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useFail,
				method: "GET",
				target: "/products/",
				body:   "",
			},
			want: want{
				code:    500,
				content: "text/plain; charset=utf-8",
				body:    "Internal Server Error\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(tt.args.use)
			fn := api.products()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.target, strings.NewReader(tt.args.body))
			fn(w, r, nil)
			CheckStatusCode(t, w, tt.want.code)
			CheckContentType(t, w, tt.want.content)
			CheckBody(t, w, tt.want.body)
		})
	}
}

func TestApiVersionsMem(t *testing.T) {
	type args struct {
		use    Usecaser
		method string
		target string
		body   string
	}
	type want struct {
		code    int
		content string
		body    string
	}
	useOk := NewUsecase(dbMem)
	useFail := usecaseFail{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/versions/",
				body:   `{"product":"p2"}`,
			},
			want: want{
				code:    200,
				content: "application/json; charset=utf-8",
				body:    `{"versions":["v21","v22"]}` + "\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/versions/",
				body:   `{"product":"p2`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/versions/",
				body:   `{}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useFail,
				method: "GET",
				target: "/versions/",
				body:   `{"product":"p2"}`,
			},
			want: want{
				code:    500,
				content: "text/plain; charset=utf-8",
				body:    "Internal Server Error\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(tt.args.use)
			fn := api.versions()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.target, strings.NewReader(tt.args.body))
			fn(w, r, nil)
			CheckStatusCode(t, w, tt.want.code)
			CheckContentType(t, w, tt.want.content)
			CheckBody(t, w, tt.want.body)
		})
	}
}

func TestApiAttributesMem(t *testing.T) {
	type args struct {
		use    Usecaser
		method string
		target string
		body   string
	}
	type want struct {
		code    int
		content string
		body    string
	}
	useOk := NewUsecase(dbMem)
	useFail := usecaseFail{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/attributes/",
				body:   `{"product":"p3","version":"v3"}`,
			},
			want: want{
				code:    200,
				content: "application/json; charset=utf-8",
				body:    `{"attributes":["a31","a32","a33"]}` + "\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/attributes/",
				body:   `{"product":"p3`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/attributes/",
				body:   `{"version":"v3"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/attributes/",
				body:   `{"product":"p3"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useFail,
				method: "GET",
				target: "/attributes/",
				body:   `{"product":"p3","version":"v3"}`,
			},
			want: want{
				code:    500,
				content: "text/plain; charset=utf-8",
				body:    "Internal Server Error\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(tt.args.use)
			fn := api.attributes()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.target, strings.NewReader(tt.args.body))
			fn(w, r, nil)
			CheckStatusCode(t, w, tt.want.code)
			CheckContentType(t, w, tt.want.content)
			CheckBody(t, w, tt.want.body)
		})
	}
}

func TestApiNamesMem(t *testing.T) {
	type args struct {
		use    Usecaser
		method string
		target string
		body   string
	}
	type want struct {
		code    int
		content string
		body    string
	}
	useOk := NewUsecase(dbMem)
	useFail := usecaseFail{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/names/",
				body:   `{"product":"p4","version":"v4","attributes":["a41","a42"]}`,
			},
			want: want{
				code:    200,
				content: "application/json; charset=utf-8",
				body:    `{"names":["n41","n42"]}` + "\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/names/",
				body:   `{"product":"p4`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/names/",
				body:   `{"version":"v4","attributes":["a41","a42"]}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/names/",
				body:   `{"product":"p4","attributes":["a41","a42"]}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useFail,
				method: "GET",
				target: "/names/",
				body:   `{"product":"p4","version":"v4","attributes":["a41","a42"]}`,
			},
			want: want{
				code:    500,
				content: "text/plain; charset=utf-8",
				body:    "Internal Server Error\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(tt.args.use)
			fn := api.names()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.target, strings.NewReader(tt.args.body))
			fn(w, r, nil)
			CheckStatusCode(t, w, tt.want.code)
			CheckContentType(t, w, tt.want.content)
			CheckBody(t, w, tt.want.body)
		})
	}
}

func TestApiMeasurementsMem(t *testing.T) {
	type args struct {
		use    Usecaser
		method string
		target string
		body   string
	}
	type want struct {
		code    int
		content string
		body    string
	}
	useOk := NewUsecase(dbMem)
	useFail := usecaseFail{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/measurements/",
				body:   `{"product":"p5","version":"v5","attributes":["a51","a52"],"name":"n5"}`,
			},
			want: want{
				code:    200,
				content: "application/json; charset=utf-8",
				body:    `{"measurements":["m51","m52"]}` + "\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/measurements/",
				body:   `{"product":"p5`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/measurements/",
				body:   `{"version":"v5","attributes":["a51","a52"],"name":"n5"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/measurements/",
				body:   `{"product":"p5","attributes":["a51","a52"],"name":"n5"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/measurements/",
				body:   `{"product":"p5","version":"v5","attributes":["a51","a52"]}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useFail,
				method: "GET",
				target: "/measurements/",
				body:   `{"product":"p5","version":"v5","attributes":["a51","a52"],"name":"n5"}`,
			},
			want: want{
				code:    500,
				content: "text/plain; charset=utf-8",
				body:    "Internal Server Error\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(tt.args.use)
			fn := api.measurements()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.target, strings.NewReader(tt.args.body))
			fn(w, r, nil)
			CheckStatusCode(t, w, tt.want.code)
			CheckContentType(t, w, tt.want.content)
			CheckBody(t, w, tt.want.body)
		})
	}
}

func TestApiChartMem(t *testing.T) {
	type args struct {
		use    Usecaser
		method string
		target string
		body   string
	}
	type want struct {
		code    int
		content string
		body    string
	}
	useOk := NewUsecase(dbMem)
	useFail := usecaseFail{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","lbound":"2017-08-19T08:40:23Z","rbound":"2018-07-02T06:34:54Z"}`,
			},
			want: want{
				code:    200,
				content: "application/json; charset=utf-8",
				body:    `{"data":[{"x":"2017-08-19T08:40:23Z","y":10},{"x":"2017-10-05T12:57:44Z","y":33},{"x":"2017-10-29T16:18:07Z","y":48}]}` + "\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","lbound":"2017-10-05T12:57:44Z","rbound":"2017-10-05T12:57:44Z"}`,
			},
			want: want{
				code:    200,
				content: "application/json; charset=utf-8",
				body:    `{"data":[{"x":"2017-10-05T12:57:44Z","y":33}]}` + "\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","lbound":"2019-03-31T11:12:16Z","rbound":"2019-04-06T14:20:19Z"}`,
			},
			want: want{
				code:    200,
				content: "application/json; charset=utf-8",
				body:    `{"data":null}` + "\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p7","version":"v7","attributes":["a71","a72"],"name":"n7","measurement":"m7","lbound":"2017-08-19T08:40:23Z","rbound":"2018-07-02T06:34:54Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Not enough attributes provided!\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","lbound":"2017-08-19T08:40:23Z","rbound":"2018-07-02T06:34:54Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","attributes":["a61","a62"],"name":"n6","measurement":"m6","lbound":"2017-08-19T08:40:23Z","rbound":"2018-07-02T06:34:54Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"measurement":"m6","lbound":"2017-08-19T08:40:23Z","rbound":"2018-07-02T06:34:54Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","lbound":"2017-08-19T08:40:23Z","rbound":"2018-07-02T06:34:54Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","rbound":"2018-07-02T06:34:54Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","lbound":"2017-08-19T08:40:23Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useOk,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","lbound":"2018-07-02T06:34:54Z","rbound":"2017-08-19T08:40:23Z"}`,
			},
			want: want{
				code:    400,
				content: "text/plain; charset=utf-8",
				body:    "Bad Request\n",
			},
		},
		{
			name: "",
			args: args{
				use:    useFail,
				method: "GET",
				target: "/chart/",
				body:   `{"product":"p6","version":"v6","attributes":["a61","a62"],"name":"n6","measurement":"m6","lbound":"2017-08-19T08:40:23Z","rbound":"2018-07-02T06:34:54Z"}`,
			},
			want: want{
				code:    500,
				content: "text/plain; charset=utf-8",
				body:    "Internal Server Error\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(tt.args.use)
			fn := api.chart()
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.target, strings.NewReader(tt.args.body))
			fn(w, r, nil)
			CheckStatusCode(t, w, tt.want.code)
			CheckContentType(t, w, tt.want.content)
			CheckBody(t, w, tt.want.body)
		})
	}
}

type usecaseFail struct{}

func (usecaseFail) Add(ctx context.Context, job *Job) error {
	return errors.New("")
}

func (usecaseFail) Products() ([]string, error) {
	return nil, errors.New("")
}

func (usecaseFail) Versions(product string) ([]string, error) {
	return nil, errors.New("")
}

func (usecaseFail) Attributes(product string, version string) ([]string, error) {
	return nil, errors.New("")
}

func (usecaseFail) Names(product string, version string, attributes []string) ([]string, error) {
	return nil, errors.New("")
}

func (usecaseFail) Measurements(product string, version string, attributes []string, name string) ([]string, error) {
	return nil, errors.New("")
}

func (usecaseFail) Jobs(product string, version string, attributes []string, name string, measurement string, lbound time.Time, rbound time.Time) ([]*Job, error) {
	return nil, errors.New("")
}
