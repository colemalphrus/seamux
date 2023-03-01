/*
seamux provides a lightweight HTTP router and middleware framework for Go.

The package defines a type RouteMux that allows adding routes and middleware to handle HTTP requests. A route is defined as a combination of a URL pattern, a set of HTTP methods, and an HTTP request handler.

To create a new RouteMux, call seamux.New() function. You can then use the HandleFunc() method to add routes to the mux.

The AddRoute() method is used internally to create a new route and add it to the list of routes in the RouteMux.

The ServeHTTP() method is the main entry point for handling HTTP requests. It matches the incoming request to the appropriate route and executes its handler. Middleware can also be executed before and after the request handler is invoked.

This implementation of HTTP routing uses regular expressions to match URL patterns and to extract parameters from URL paths. Parameters are defined as segments of the URL path that start with a colon (e.g. "/users/:id"). These parameters are then mapped to their corresponding values and made available to the request handler through the request object.
*/

package seamux

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type route struct {
	methods []string
	regex   *regexp.Regexp
	params  map[int]string
	handler http.HandlerFunc
}

type RouteMux struct {
	routes     []*route
	middleware []http.HandlerFunc
}

func New() *RouteMux {
	return &RouteMux{}
}

func (m *RouteMux) HandleFunc(pattern string, handler http.HandlerFunc) *route {
	return m.AddRoute(pattern, handler)
}

func (r *route) Methods(methods ...string) {
	r.methods = methods
}

func (m *RouteMux) Middleware(filter http.HandlerFunc) {
	m.middleware = append(m.middleware, filter)
}

func (m *RouteMux) AddRoute(pattern string, handler http.HandlerFunc) *route {

	//split the url into sections
	parts := strings.Split(pattern, "/")

	//find params that start with ":"
	//replace with regular expressions
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			params[j] = part[1:]
			parts[i] = "([^/]+)"
			j++
		}
	}

	//recreate the url pattern, with parameters replaced
	//by regular expressions. then compile the regex
	pattern = strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
	}

	//now create the Route
	route := &route{}
	route.regex = regex
	route.handler = handler
	route.params = params

	//and finally append to the list of Routes
	m.routes = append(m.routes, route)

	return route
}

func (m *RouteMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestPath := r.URL.Path

	//find a matching Route
	for _, route := range m.routes {

		if !validateMethod(route.methods, r.Method) {
			continue
		}

		if !route.regex.MatchString(requestPath) {
			continue
		}

		//get path params
		matches := route.regex.FindStringSubmatch(requestPath)

		if len(route.params) > 0 {
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(route.params[i], match)
			}

			r.URL.RawQuery = url.Values(values).Encode()
		}

		//execute middleware
		for _, filter := range m.middleware {
			filter(w, r)
		}

		//Invoke the request handler
		route.handler(w, r)
		break
	}
}

func validateMethod(s []string, str string) bool {
	if len(s) == 0 {
		return true
	}
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
