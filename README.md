# SeaMux
### 🚢  SeaMUX: The Multiplexing Mate for Your HTTP Voyage! 🚢

Package SeaMUX provides a lightweight and flexible HTTP router for Go. It is designed to be simple to use, yet powerful enough to handle complex routing requirements.

## Usage

### New() *RouteMux
Creates and returns a new RouteMux.
```go
mux := seamux.New()
```

### HandleFunc(pattern string, handler http.HandlerFunc) *route
Adds a new route to the RouteMux.
```go
mux.HandleFunc("/hello/:name", func(w http.ResponseWriter, r *http.Request) {
name := r.URL.Query().Get("name")
fmt.Fprintf(w, "Hello, %s!", name)
})
```
- pattern : the URL pattern to match. Any sections that start with : are considered parameters and are replaced with a regular expression.
- handler : the function to execute when the pattern is matched.
- Returns a pointer to the route object that was created.

### Middleware(filter http.HandlerFunc)
Adds a middleware function to the RouteMux.

```go
mux.Middleware(func(w http.ResponseWriter, r *http.Request) {
// Do something before the request is handled by the route's handler.
})
```
- filter : the middleware function to add.

### AddRoute(pattern string, handler http.HandlerFunc) *route
Adds a new route to the RouteMux using the specified pattern and handler. This method is called by HandleFunc.

```go
route := mux.AddRoute("/hello/:name", func(w http.ResponseWriter, r *http.Request) {
name := r.URL.Query().Get("name")
fmt.Fprintf(w, "Hello, %s!", name)
})
```
- pattern : the URL pattern to match. Any sections that start with : are considered parameters and are replaced with a regular expression.
- handler : the function to execute when the pattern is matched.
- Returns a pointer to the route object that was created.

### ServeHTTP(w http.ResponseWriter, r *http.Request)
Handles incoming HTTP requests and routes them to the appropriate handler.
```go
http.ListenAndServe(":8080", mux)
```
- w : the http.ResponseWriter to use for the response.
- r : the incoming http.Request to route.


### Types
####  route
```go
type route struct {
methods []string
regex   *regexp.Regexp
params  map[int]string
handler http.HandlerFunc
}
```
The route type represents an HTTP route, consisting of a set of HTTP methods, a regular expression for matching the URL pattern, a set of named parameters extracted from the URL pattern, and an HTTP request handler function.

#### RouteMux
```go

type RouteMux struct {
routes     []*route
middleware []http.HandlerFunc
}
```
The RouteMux type is the main HTTP multiplexer type. It stores a list of routes and middleware functions.

### Functions

```go
func New() *RouteMux
```
New returns a new RouteMux.

```go
func (m *RouteMux) HandleFunc
```
HandleFunc mimics the HandleFunc in the standard http library

```go
func (m *RouteMux) HandleFunc(pattern string, handler http.HandlerFunc) *route
```
HandleFunc registers an HTTP request handler function for the given URL pattern. It returns a route that can be used to further configure the route.

```go
func (r *route) Methods(methods ...string)
```
Methods sets the HTTP methods that this route should match.

```go
func (m *RouteMux) Middleware(filter http.HandlerFunc)
```
Middleware registers an HTTP middleware function to be executed for all routes.

```go
func (m *RouteMux) AddRoute(pattern string, handler http.HandlerFunc) *route
```
AddRoute registers a new HTTP route for the given URL pattern and request handler function. It returns a route that can be used to further configure the route.

```go
func (m *RouteMux) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP is the HTTP request handler function for the RouteMux type. It matches incoming HTTP requests to registered routes, executes middleware functions, and invokes the appropriate request handler function.

```go
func containsMethod(s []string, str string) bool
```
containsMethod returns true if the given HTTP method is in the list of allowed methods for a route, false otherwise.