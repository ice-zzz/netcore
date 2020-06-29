/**
 *     ______                 __
 *    /\__  _\               /\ \
 *    \/_/\ \/     ___     __\ \ \         __      ___     ___     __
 *       \ \ \    / ___\ / __ \ \ \  __  / __ \  /  _  \  / ___\ / __ \
 *        \_\ \__/\ \__//\  __/\ \ \_\ \/\ \_\ \_/\ \/\ \/\ \__//\  __/
 *        /\_____\ \____\ \____\\ \____/\ \__/ \_\ \_\ \_\ \____\ \____\
 *        \/_____/\/____/\/____/ \/___/  \/__/\/_/\/_/\/_/\/____/\/____/
 *
 *
 *                                                                    @寒冰
 *                                                            www.icezzz.cn
 *                                                     hanbin020706@163.com
 */
package rest

type AllowOriginMiddleware struct {
	Origin  string
	Methods string
	Headers string
}

func (ao *AllowOriginMiddleware) MiddlewareFunc(h HandlerFunc) HandlerFunc {
	if ao.Origin == "" {
		ao.Origin = "*"
	}
	if ao.Methods == "" {
		ao.Methods = "POST, GET, PUT, DELETE"
	}
	if ao.Headers == "" {
		ao.Headers = "Action, Module"
	}

	return func(w ResponseWriter, r *Request) {

		w.Header().Set("Access-Control-Allow-Origin", ao.Origin)
		w.Header().Set("Access-Control-Allow-Methods", ao.Methods)
		w.Header().Set("Access-Control-Allow-Headers", ao.Headers)

		// call the handler
		h(w, r)

	}
}
