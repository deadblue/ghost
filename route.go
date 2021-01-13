package ghost

import (
	"github.com/deadblue/ghost/internal/context"
	"github.com/deadblue/ghost/internal/rule"
	"log"
)

type _RouteKey struct {
	method, path string
}

type _RoutePath struct {
	// Path controller
	ctrl Controller
	// Path rule
	rule *rule.Rule
}

type _RouteTable struct {
	// Exact path mapping
	mapping map[_RouteKey]Controller
	// Branches
	branches map[string][]*_RoutePath
}

// Mount mounts controller into a request path which is described by name.
func (t *_RouteTable) Mount(name string, ctrl Controller) (err error) {
	// Parse rule
	m, r, err := rule.Parse(name)
	if err != nil {
		log.Printf("Parse rule error: %s", err)
		return
	}
	// Store exactly matching in map
	if path, exactly := r.Path(); exactly {
		t.mapping[_RouteKey{m, path}] = ctrl
		return
	}
	// In other case, store it in branches
	branch, exists := t.branches[m]
	if !exists {
		branch = make([]*_RoutePath, 0)
	}
	branch = append(branch, &_RoutePath{
		ctrl: ctrl,
		rule: r,
	})
	t.branches[m] = branch
	return
}

func (t *_RouteTable) Resolve(ctx *context.Impl) (ctrl Controller) {
	// Get request method and path
	rm, rp := ctx.MethodAndPath()
	// First search exactly matching
	exists := false
	if ctrl, exists = t.mapping[_RouteKey{rm, rp}]; exists {
		return ctrl
	}
	// Search branch
	branch, exists := t.branches[rm]
	if !exists {
		return nil
	}
	// Scan request path one by one.
	maxScore, pathVars := -1, map[string]string{}
	for _, path := range branch {
		vars := make(map[string]string)
		score := path.rule.Match(rp, vars)
		// Skip mismatching and low-quality matching
		if score >= 0 && score > maxScore {
			maxScore, pathVars, ctrl = score, vars, path.ctrl
		}
	}
	if maxScore >= 0 {
		for k, v := range pathVars {
			ctx.PutPathVar(k, v)
		}
	}
	return
}
