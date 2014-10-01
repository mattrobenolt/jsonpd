package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"sync/atomic"
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

var (
	COUNTER_OK,
	COUNTER_E_EMPTY,
	COUNTER_E_TOOLONG,
	COUNTER_E_RESERVED,
	COUNTER_E_INVALID uint64
	BEGIN = time.Now()
)

type Handler struct{}

func (*Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// pop off key
	cb := q.Get(*callback)
	q.Del(*callback)

	if err := isValid(cb); err != nil {
		switch err {
		case E_EMPTY:
			atomic.AddUint64(&COUNTER_E_EMPTY, 1)
		case E_TOOLONG:
			atomic.AddUint64(&COUNTER_E_TOOLONG, 1)
			cb = cb[:50]
		case E_RESERVED:
			atomic.AddUint64(&COUNTER_E_RESERVED, 1)
		case E_INVALID:
			atomic.AddUint64(&COUNTER_E_INVALID, 1)
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

	atomic.AddUint64(&COUNTER_OK, 1)
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

func writeStats(w io.Writer) {
	w.Write([]byte(fmt.Sprintln("uptime:", int64(time.Now().Sub(BEGIN).Seconds()))))
	w.Write([]byte(fmt.Sprintln("ok:", COUNTER_OK)))
	w.Write([]byte(fmt.Sprintln("e_empty:", COUNTER_E_EMPTY)))
	w.Write([]byte(fmt.Sprintln("e_toolong:", COUNTER_E_TOOLONG)))
	w.Write([]byte(fmt.Sprintln("e_reserved:", COUNTER_E_RESERVED)))
	w.Write([]byte(fmt.Sprintln("e_invalid:", COUNTER_E_INVALID)))
}

var (
	bind     = flag.String("b", "localhost:8000", "bind address (default: localhost:8000)")
	procs    = flag.Int("n", runtime.NumCPU(), fmt.Sprintf("num procs (default: %d)", runtime.NumCPU()))
	callback = flag.String("cb", "callback", "callback argument (default: callback)")
	timeout  = flag.Int("t", 500, "timeout (default: 500)")
	info     = flag.String("i", "localhost:8001", "bind address for stats (default: localhost:8001)")
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
	fmt.Println("info:", *info)
	fmt.Println("callback:", *callback)
	fmt.Println("timeout:", time.Duration(*timeout)*time.Millisecond)

	s := &http.Server{
		Addr:         *bind,
		Handler:      &Handler{},
		ReadTimeout:  time.Duration(*timeout) * time.Millisecond,
		WriteTimeout: time.Duration(*timeout) * time.Millisecond,
	}

	ln, err := net.Listen("tcp", *info)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			writeStats(conn)
			conn.Close()
		}
	}()

	fmt.Println("")
	log.Println("ready.")
	log.Fatal(s.ListenAndServe())
}
