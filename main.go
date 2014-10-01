package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"time"
)

var (
	RESERVED_WORDS = map[string]struct{}{
		"abstract":     struct{}{},
		"boolean":      struct{}{},
		"break":        struct{}{},
		"byte":         struct{}{},
		"case":         struct{}{},
		"catch":        struct{}{},
		"char":         struct{}{},
		"class":        struct{}{},
		"const":        struct{}{},
		"continue":     struct{}{},
		"debugger":     struct{}{},
		"default":      struct{}{},
		"delete":       struct{}{},
		"do":           struct{}{},
		"double":       struct{}{},
		"else":         struct{}{},
		"enum":         struct{}{},
		"export":       struct{}{},
		"extends":      struct{}{},
		"false":        struct{}{},
		"final":        struct{}{},
		"finally":      struct{}{},
		"float":        struct{}{},
		"for":          struct{}{},
		"function":     struct{}{},
		"goto":         struct{}{},
		"if":           struct{}{},
		"implements":   struct{}{},
		"import":       struct{}{},
		"in":           struct{}{},
		"instanceof":   struct{}{},
		"int":          struct{}{},
		"interface":    struct{}{},
		"long":         struct{}{},
		"native":       struct{}{},
		"new":          struct{}{},
		"null":         struct{}{},
		"package":      struct{}{},
		"private":      struct{}{},
		"protected":    struct{}{},
		"public":       struct{}{},
		"return":       struct{}{},
		"short":        struct{}{},
		"static":       struct{}{},
		"super":        struct{}{},
		"switch":       struct{}{},
		"synchronized": struct{}{},
		"this":         struct{}{},
		"throw":        struct{}{},
		"throws":       struct{}{},
		"transient":    struct{}{},
		"true":         struct{}{},
		"try":          struct{}{},
		"typeof":       struct{}{},
		"var":          struct{}{},
		"void":         struct{}{},
		"volatile":     struct{}{},
		"while":        struct{}{},
		"with":         struct{}{},
	}

	CALLBACK_RE = regexp.MustCompile(`^[$a-zA-Z_][0-9a-zA-Z_\.\[\]]*$`)
)

type Handler struct{}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// pop off key
	cb := q.Get(*callback)
	q.Del(*callback)

	if err := isValid(cb); err != nil {
		if err == E_TOOLONG {
			cb = cb[:50]
		}
		log.Printf("[reject] %s: %s", err, cb)
		handle400(w, r)
		return
	}

	qs := q.Encode()

	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Disposition", "attachment; filename=f.txt")
	w.Write([]byte("<!--esi\n/**/"))
	w.Write([]byte(cb))
	w.Write([]byte("(<esi:include src=\""))
	w.Write([]byte(r.URL.Path))
	if qs != "" {
		w.Write([]byte{'?'})
		w.Write([]byte(qs))
	}
	w.Write([]byte("\"/>);\n-->\n"))
}

var (
	E_EMPTY    = errors.New("empty")
	E_TOOLONG  = errors.New("too long")
	E_RESERVED = errors.New("reserved")
	E_INVALID  = errors.New("invalid")
)

func isValid(cb string) error {
	if cb == "" {
		return E_EMPTY
	}
	// Callbacks longer than 50 characters are suspicious.
	// There isn't a legit reason for a callback longer.
	// The length is arbitrary too.
	// It's technically possible to construct malicious payloads using
	// only ascii characters, so we just block this.
	if len(cb) > 50 {
		return E_TOOLONG
	}

	if _, ok := RESERVED_WORDS[cb]; ok {
		return E_RESERVED
	}

	if ok := CALLBACK_RE.MatchString(cb); !ok {
		return E_INVALID
	}

	return nil
}

func handle400(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Connection", "close")
	w.WriteHeader(400)
	w.Write([]byte("<html><body><h1>400 Bad Request</h1></body></html>\n"))
}

var (
	bind     = flag.String("b", ":8000", "bind address (default: :8000)")
	procs    = flag.Int("n", runtime.NumCPU(), fmt.Sprintf("num procs (default: %d)", runtime.NumCPU()))
	callback = flag.String("cb", "callback", "callback argument (default: callback)")
	timeout  = flag.Int("t", 500, "timeout (default: 500)")
)

func init() {
	flag.Parse()
}

func art() {
	fmt.Println(`
     ██╗███████╗ ██████╗ ███╗   ██╗██████╗ ██████╗
     ██║██╔════╝██╔═══██╗████╗  ██║██╔══██╗██╔══██╗
     ██║███████╗██║   ██║██╔██╗ ██║██████╔╝██║  ██║
██   ██║╚════██║██║   ██║██║╚██╗██║██╔═══╝ ██║  ██║
╚█████╔╝███████║╚██████╔╝██║ ╚████║██║     ██████╔╝
 ╚════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝     ╚═════╝
                                                   `)
}

func main() {
	art()
	runtime.GOMAXPROCS(*procs)

	fmt.Println("procs:", *procs)
	fmt.Println("bind:", *bind)
	fmt.Println("callback:", *callback)
	fmt.Println("timeout:", time.Duration(*timeout)*time.Millisecond)

	s := &http.Server{
		Addr:         *bind,
		Handler:      &Handler{},
		ReadTimeout:  time.Duration(*timeout) * time.Millisecond,
		WriteTimeout: time.Duration(*timeout) * time.Millisecond,
	}
	fmt.Println("")
	log.Println("ready.")
	log.Fatal(s.ListenAndServe())
}
