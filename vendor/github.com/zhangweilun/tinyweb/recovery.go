package tinyweb

import (
	"log"
	"os"
	"net/http"
	"runtime"
	"fmt"
)

/**
* 
* @author willian
* @created 2017-01-27 22:34
* @email 18702515157@163.com  
**/
// Recovery is a Tinyweb middleware that recovers from any panics and writes a 500 if there was one.

type Recovery struct {
	Logger *log.Logger
	PrintStack bool
	StackAll bool
	StackSize int
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "[negroni] ", 0),
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "PANIC: %s\n%s"
			rec.Logger.Printf(f, err, stack)

			if rec.PrintStack {
				fmt.Fprintf(rw, f, err, stack)
			}
		}
	}()

	next(rw, r)
}