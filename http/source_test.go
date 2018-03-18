package http_test

import (
	"fmt"
	"net"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/pilosa/pdk/http"
)

func TestJSONSource(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listening: %v", err)
	}

	j, err := http.NewJSONSource(http.WithListener(ln))
	if err != nil {
		t.Fatalf("getting json source: %v", err)
	}

	tests := []struct {
		method string
		path   string
		data   string
		exp    map[string]interface{}
		expErr string
	}{
		{
			method: "POST",
			path:   "/",
			data:   `{"hello": 2}`,
			exp:    map[string]interface{}{"hello": 2.0},
			expErr: "",
		},
		{
			method: "POST",
			path:   "/blah",
			data:   `{"hello": 2}`,
			exp:    map[string]interface{}{"hello": 2.0},
			expErr: "",
		},
		// // TODO we now ignore errors in the JSONSource, so these test are
		// // commented. need to decide on an actual strategy for error handling
		// // and reporting.
		// {
		// 	method: "POST",
		// 	path:   "/",
		// 	data:   `{"hello: 2}`,
		// 	exp:    nil,
		// 	expErr: "decoding json",
		// },
		// {
		// 	method: "GET",
		// 	path:   "/",
		// 	data:   ``,
		// 	exp:    nil,
		// 	expErr: "unsupported method",
		// },
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			j.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(test.method, test.path, strings.NewReader(test.data)))
			data, err := j.Record()
			if err != nil {
				if test.expErr == "" {
					t.Fatalf("unexpected error: %v", err)
				}
				if !strings.Contains(err.Error(), test.expErr) {
					t.Fatalf("wrong err: %v", err)
				}
			} else if !reflect.DeepEqual(data, test.exp) {
				t.Fatalf("unexpected data: %#v", data)
			} else if err == nil && test.expErr != "" {
				t.Fatalf("did not get expected err")
			}
		})
	}

}