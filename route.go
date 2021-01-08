package ghost

import (
	"github.com/deadblue/ghost/internal/rule"
	"log"
	"net/http"
	"strings"
)

type _RoutePath struct {
	// Path controller
	c Controller
	// Path rule
	r *rule.Rule
}

type _RouteBranch struct {
	// Root controller
	c Controller
	// Sub path
	ps []*_RoutePath
}

type _RouteTree struct {
	// Branches
	bs map[string]*_RouteBranch
}

func (t *_RouteTree) Mount(name string, ctrl Controller) (err error) {
	// Parse rule
	m, r, err := rule.Parse(name)
	if err != nil {
		log.Printf("Parse rule error: %s", err)
		return
	}
	// Mount to tree
	if t.bs == nil {
		t.bs = map[string]*_RouteBranch{}
	}
	branch, exists := t.bs[m]
	if !exists {
		branch = &_RouteBranch{
			ps: []*_RoutePath{},
		}
		t.bs[m] = branch
	}
	if r == nil {
		branch.c = ctrl
	} else {
		branch.ps = append(branch.ps, &_RoutePath{
			c: ctrl,
			r: r,
		})
	}
	return
}

func (t *_RouteTree) Resolve(r *http.Request, ctx *_ContextImpl) (ctrl Controller) {
	// Get request method and path
	rm, rp := strings.ToLower(r.Method), r.URL.Path
	branch, exists := t.bs[rm]
	if !exists {
		return nil
	}
	if rp == "/" {
		return branch.c
	}
	// Match request path one by one.
	// TODO:
	//  * Consider add cache for matched path.
	//  * Optimize the algorithm, to make it faster than O(n).
	// Since the controllers are not so many, O(n) is enough now.
	maxScore, pathVars := 0, map[string]string{}
	for _, p := range branch.ps {
		vars := make(map[string]string)
		score := p.r.Match(rp, vars)
		if score != 0 && score > maxScore {
			ctrl = p.c
			maxScore, pathVars = score, vars
		}
	}
	if maxScore > 0 {
		for k, v := range pathVars {
			ctx.pv[k] = v
		}
	}
	return
}
