package rest

import (
	"fmt"
	"net/http"

	"github.com/ice-zzz/netcore/service"
)

// Api defines a stack of Middlewares and an App.
type Api struct {
	stack []Middleware
	app   App
	serv  *http.Server
	service.Entity
}

// NewApi makes a new Api object. The Middleware stack is empty, and the App is nil.
func NewApi() *Api {
	return &Api{
		stack: []Middleware{},
		app:   nil,
		serv:  nil,
		Entity: service.Entity{
			Name: "",
			Ip:   "0.0.0.0",
			Port: 5679,
		},
	}
}

func (api *Api) Start() {
	api.serv = &http.Server{
		Addr: fmt.Sprintf("%s:%d", api.Ip, api.Port),
	}
	_ = api.serv.ListenAndServe()
}

func (api *Api) Stop() {
	_ = api.serv.Close()
}

// Use pushes one or multiple middlewares to the stack for middlewares
// maintained in the Api object.
func (api *Api) Use(middlewares ...Middleware) {
	api.stack = append(api.stack, middlewares...)
}

// SetApp sets the App in the Api object.
func (api *Api) SetApp(app App) {
	api.app = app
}

// MakeHandler wraps all the Middlewares of the stack and the App together, and returns an
// http.Handler ready to be used. If the Middleware stack is empty the App is used directly. If the
// App is nil, a HandlerFunc that does nothing is used instead.
func (api *Api) MakeHandler() http.Handler {
	var appFunc HandlerFunc
	if api.app != nil {
		appFunc = api.app.AppFunc()
	} else {
		appFunc = func(w ResponseWriter, r *Request) {}
	}
	return http.HandlerFunc(
		adapterFunc(
			WrapMiddlewares(api.stack, appFunc),
		),
	)
}

// Defines a stack of middlewares convenient for development. Among other things:
// console friendly logging, JSON indentation, error stack strace in the response.
var DefaultDevStack = []Middleware{
	&AccessLogApacheMiddleware{},
	&TimerMiddleware{},
	&RecorderMiddleware{},
	&PoweredByMiddleware{},
	&RecoverMiddleware{
		EnableResponseStackTrace: true,
	},
	&JsonIndentMiddleware{},
	&ContentTypeCheckerMiddleware{},
}

// Defines a stack of middlewares convenient for production. Among other things:
// Apache CombinedLogFormat logging, gzip compression.
var DefaultProdStack = []Middleware{
	&AccessLogApacheMiddleware{
		Format: CombinedLogFormat,
	},
	&TimerMiddleware{},
	&RecorderMiddleware{},
	&PoweredByMiddleware{},
	&RecoverMiddleware{},
	&GzipMiddleware{},
	&ContentTypeCheckerMiddleware{},
}

// Defines a stack of middlewares that should be common to most of the middleware stacks.
var DefaultCommonStack = []Middleware{
	&TimerMiddleware{},
	&RecorderMiddleware{},
	&PoweredByMiddleware{},
	&RecoverMiddleware{},
}
