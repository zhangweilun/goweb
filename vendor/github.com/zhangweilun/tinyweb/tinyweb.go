package tinyweb

import (
	"net/http"
	"log"
	"os"
)

/**
* 
* @author willian
* @created 2017-01-27 19:03
* @email 18702515157@163.com  
**/

type Handler interface {
	ServeHTTP(rw http.ResponseWriter,r *http.Request,next http.HandlerFunc)
}

// HandlerFunc implements the Handler
type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter,r *http.Request,next http.HandlerFunc) {
	h(rw,r,next)
}

type middleware struct {
	handler Handler
	next *middleware
}

func (m middleware) ServeHTTP(rw http.ResponseWriter,r *http.Request)  {
	m.handler.ServeHTTP(rw,r,m.next.ServeHTTP)
}

// Wrap converts a http.Handler into a negroni.Handler so it can be used as a tinyweb
// middleware. The next http.HandlerFunc is automatically called after the Handler
// is executed.
func Wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

type Tinyweb struct {
	middleware middleware
	handlers []Handler
}


// New returns a new Negroni instance with no middleware preconfigured.
func New(handlers ...Handler) *Tinyweb {
	return &Tinyweb{
		handlers:   handlers,
		middleware: build(handlers),
	}
}

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Classic() *Tinyweb {
	return New(NewRecovery(), NewLogger(), NewStatic(http.Dir("public")))
}

func (t *Tinyweb) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	t.middleware.ServeHTTP(NewResponseWriter(rw),r)
}


func (t *Tinyweb) Add(h Handler) {
	if h == nil{
		panic("handler cannot be nil")
	}else {
		t.handlers = append(t.handlers,h)
		t.middleware = build(t.handlers)
	}
}

// UseHandler adds a http.Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (t *Tinyweb) UseHandler(handler http.Handler) {
	t.Add(Wrap(handler))
}

func (t *Tinyweb) Run(port string)  {
	l := log.New(os.Stdout, "[tinyweb] ", 0)
	l.Printf("listening on %s", port)
	l.Fatal(http.ListenAndServe(port, t))
}



func build(handlers []Handler) middleware {
	var next middleware

	if len(handlers) == 0 {
		return voidMiddleware()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = voidMiddleware()
	}

	return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
	return middleware{
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
		&middleware{},
	}
}