package engine

type (
	// Group is a set of sub-routes for a specified route. It can be used for inner
	// routes that share a common middleware or functionality that should be separate
	// from the parent engine instance while still inheriting from it.
	Group struct {
		common
		prefix     string
		middleware []MiddlewareFunc
		engine     *Engine
	}
)

func newGroup(prefix string, e *Engine) *Group {
	g := &Group{prefix: prefix, engine: e}
	g.common.add = g.Add
	return g
}

// Use implements `Engine#Use()` for sub-routes within the Group.
func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
	if len(g.middleware) == 0 {
		return
	}
	// Allow all requests to reach the group as they might get dropped if router
	// doesn't find a match, making none of the group middleware process.
	g.Any("", NotFoundHandler)
	g.Any("/*", NotFoundHandler)
}

// Add implements `Engine#Add()` for sub-routes within the Group.
func (g *Group) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	m := make([]MiddlewareFunc, 0, len(g.middleware)+len(middleware))
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	g.engine.Add(method, g.prefix+path, handler, m...)
}