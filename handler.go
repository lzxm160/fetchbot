package fetchbot

import (
	"net/http"
	"strings"
	"sync"
)

// Context is a Command's fetch context, passed to the Handler. It gives access to the
// original Command and the associated Queue.
type Context struct {
	Cmd Command
	Q   *Queue
}

// The Handler interface is used to process the Fetcher's requests. It is similar to the
// net/http.Handler interface.
type Handler interface {
	Handle(*Context, *http.Response, error)
}

// A HandlerFunc is a function signature that implements the Handler interface. A function
// with this signature can thus be used as a Handler.
type HandlerFunc func(*Context, *http.Response, error)

// Handle is the Handler interface implementation for the HandlerFunc type.
func (h HandlerFunc) Handle(ctx *Context, res *http.Response, err error) {
	h(ctx, res, err)
}

// Mux is a simple multiplexer for the Handler interface, similar to net/http.ServeMux.
// It is itself a Handler, and dispatches the calls to the matching Handlers.
type Mux struct {
	DefaultHandler Handler

	mu   sync.RWMutex
	errm map[error]Handler
	res  map[*resEntry]bool // a set of entries
}

// NewMux returns an initialized Mux.
func NewMux() *Mux {
	return &Mux{
		// Default handler is a no-op
		DefaultHandler: HandlerFunc(func(ctx *Context, res *http.Response, err error) {}),
		errm:           make(map[error]Handler),
		res:            make(map[*resEntry]bool),
	}
}

// Handle is the Handler interface implementation for Mux. It dispatches the calls
// to the matching Handler.
func (mux *Mux) Handle(ctx *Context, res *http.Response, err error) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	if err != nil {
		// Find a matching error handler
		if h, ok := mux.errm[err]; ok {
			h.Handle(ctx, res, err)
			return
		}
		if h, ok := mux.errm[nil]; ok {
			h.Handle(ctx, res, err)
			return
		}
	} else {
		// Find a matching response handler
		var h Handler
		var n = -1
		for r, _ := range mux.res {
			if ok, cnt := r.match(res); ok {
				if cnt > n {
					h, n = r.h, cnt
				}
			}
		}
		if h != nil {
			h.Handle(ctx, res, err)
			return
		}
	}
	mux.DefaultHandler.Handle(ctx, res, err)
}

// HandleError registers a Handler for a specific error value. Multiple calls
// with the same error value override previous calls. As a special case, a nil
// error value registers a Handler for any error that doesn't have a specific
// Handler.
func (mux *Mux) HandleError(err error, h Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	mux.errm[err] = h
}

// HandleErrors registers a Handler for any error that doesn't have a specific
// Handler.
func (mux *Mux) HandleErrors(h Handler) {
	mux.HandleError(nil, h)
}

// Response initializes an entry for a Response Handler based on various criteria.
// The Response Handler is not registered until Handle is called.
func (mux *Mux) Response() *resEntry {
	return &resEntry{mux: mux}
}

// resEntry holds the criteria for the response Handler.
type resEntry struct {
	method      string
	contentType string
	minStatus   int
	maxStatus   int
	path        string
	host        string
	h           Handler
	mux         *Mux
}

// match indicates if the response Handler matches the provided response, and if so,
// and if a path criteria is specified, it also indicates the length of the path match.
func (r *resEntry) match(res *http.Response) (bool, int) {
	if r.method != "" {
		if r.method != res.Request.Method {
			return false, 0
		}
	}
	if r.contentType != "" {
		if r.contentType != getContentType(res.Header.Get("Content-Type")) {
			return false, 0
		}
	}
	if r.minStatus != 0 || r.maxStatus != 0 {
		if res.StatusCode < r.minStatus || res.StatusCode > r.maxStatus {
			return false, 0
		}
	}
	if r.host != "" {
		if res.Request.URL.Host != r.host {
			return false, 0
		}
	}
	if r.path != "" {
		if strings.HasPrefix(res.Request.URL.Path, r.path) {
			return true, len(r.path)
		} else {
			return false, 0
		}
	}
	return true, 0
}

// Returns the content type stripped of any additional parameters (following the ;).
func getContentType(val string) string {
	args := strings.SplitN(val, ";", 2)
	if len(args) > 0 {
		return strings.TrimSpace(args[0])
	}
	return val
}

// Method sets a method criteria for the Response Handler. Its Handler will only be called
// if it has this HTTP method (i.e. "GET", "HEAD", ...).
func (r *resEntry) Method(m string) *resEntry {
	r.mux.mu.Lock()
	defer r.mux.mu.Unlock()
	r.method = m
	return r
}

// ContentType sets a criteria based on the Content-Type header for the Response Handler.
// Its Handler will only be called if it has this content type, ignoring any additional
// parameter on the Header value (following the semicolon, i.e. "text/html; charset=utf-8").
func (r *resEntry) ContentType(ct string) *resEntry {
	r.mux.mu.Lock()
	defer r.mux.mu.Unlock()
	r.contentType = ct
	return r
}

// Status sets a criteria based on the Status code of the response for the Response Handler.
// Its Handler will only be called if the response has this status code.
func (r *resEntry) Status(code int) *resEntry {
	r.mux.mu.Lock()
	defer r.mux.mu.Unlock()
	r.minStatus = code
	r.maxStatus = code
	return r
}

// StatusRange sets a criteria based on the Status code of the response for the Response Handler.
// Its Handler will only be called if the response has a status code between the min and max.
// If min is greater than max, the values are switched.
func (r *resEntry) StatusRange(min, max int) *resEntry {
	if min > max {
		min, max = max, min
	}
	r.mux.mu.Lock()
	defer r.mux.mu.Unlock()
	r.minStatus = min
	r.maxStatus = max
	return r
}

// Path sets a criteria based on the path of the URL for the Response Handler. Its Handler
// will only be called if the path of the URL starts with this path. Longer matches
// have priority over shorter ones.
func (r *resEntry) Path(p string) *resEntry {
	r.mux.mu.Lock()
	defer r.mux.mu.Unlock()
	r.path = p
	return r
}

// Host sets a criteria based on the host of the URL for the Response Handler. Its Handler
// will only be called if the host of the URL matches exactly the specified host.
func (r *resEntry) Host(host string) *resEntry {
	r.mux.mu.Lock()
	defer r.mux.mu.Unlock()
	r.host = host
	return r
}

// Handler sets the Handler to be called when this Response Handler is the match for
// a given response. It registers the Response Handler in its parent Mux.
func (r *resEntry) Handler(h Handler) *resEntry {
	r.mux.mu.Lock()
	defer r.mux.mu.Unlock()
	r.h = h
	if !r.mux.res[r] {
		r.mux.res[r] = true
	}
	return r
}

// HandlerFunc sets the Handler function to be called when this Response Handler is the
// match for a given response. It registers the Response Handler in its parent Mux.
func (r *resEntry) HandlerFunc(fn func(*Context, *http.Response, error)) *resEntry {
	return r.Handler(HandlerFunc(fn))
}