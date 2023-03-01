package seamux

import (
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *RouteMux
	}{
		{
			name: "Happy Path",
			want: &RouteMux{
				routes:     nil,
				middleware: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouteMux_AddRoute(t *testing.T) {
	aliveRegex, _ := regexp.Compile("/alive")
	complexRegex, _ := regexp.Compile("/alive/([^/]+)/test")
	type fields struct {
		routes     []*route
		middleware []http.HandlerFunc
	}
	type args struct {
		pattern string
		handler http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *route
	}{
		{
			name: "happy path",
			fields: fields{
				routes:     nil,
				middleware: nil,
			},
			args: args{
				pattern: "/alive",
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				},
			},
			want: &route{
				methods: nil,
				regex:   aliveRegex,
				params:  nil,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				},
			},
		},
		{
			name: "complex regex path",
			fields: fields{
				routes:     nil,
				middleware: nil,
			},
			args: args{
				pattern: "/alive/:id/test",
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				},
			},
			want: &route{
				methods: nil,
				regex:   complexRegex,
				params:  nil,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &RouteMux{
				routes:     tt.fields.routes,
				middleware: tt.fields.middleware,
			}
			got := m.AddRoute(tt.args.pattern, tt.args.handler)

			for i, _ := range got.methods {
				if got.methods[i] != tt.want.methods[i] {
					t.Errorf("methods = |||%v|||, want |||%v|||", got.methods[i], tt.want.methods[i])
				}
			}

			if got.regex.String() != tt.want.regex.String() {
				t.Errorf("regex = |||%v|||, want |||%v|||", got.regex, tt.want.regex)
			}
		})
	}
}

func TestRouteMux_HandleFunc(t *testing.T) {
	type fields struct {
		routes     []*route
		middleware []http.HandlerFunc
	}
	type args struct {
		pattern string
		handler http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "happy path",
			fields: fields{
				routes:     nil,
				middleware: nil,
			},
			args: args{
				pattern: "/alive",
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(200)
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &RouteMux{
				routes:     tt.fields.routes,
				middleware: tt.fields.middleware,
			}
			m.HandleFunc(tt.args.pattern, tt.args.handler)

			if len(m.routes) != tt.want {
				t.Errorf("HandleFunc() = %v, want %v", len(m.routes), tt.want)
			}
		})
	}
}

func Test_validateMethod(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string is in methods",
			args: args{
				s:   []string{"GET", "POST"},
				str: "GET",
			},
			want: true,
		},
		{
			name: "methods array is empty",
			args: args{
				s:   []string{},
				str: "GET",
			},
			want: true,
		},
		{
			name: "methods is not in array",
			args: args{
				s:   []string{"POST"},
				str: "GET",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateMethod(tt.args.s, tt.args.str); got != tt.want {
				t.Errorf("validateMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
