package rule

import (
	"fmt"
	"strings"
)

type Segment struct {
	IsVar bool
	Name  string
	Ext   string
}

func (s *Segment) IsValid() bool {
	return s.Name != "" || s.Ext != ""
}

func (s *Segment) String() string {
	buf := strings.Builder{}
	if s.IsVar {
		buf.WriteRune('{')
	}
	_, _ = buf.WriteString(s.Name)
	if s.IsVar {
		buf.WriteRune('}')
	}
	if s.Ext != "" {
		_, _ = buf.WriteRune('.')
		_, _ = buf.WriteString(s.Ext)
	}
	return buf.String()
}

type Rule struct {
	// Static flag
	IsStatic bool
	// Request method
	Method string
	// Request path
	Path string
	// Request path segments
	Segments []*Segment
}

func (r *Rule) StaticKey() string {
	return fmt.Sprintf("%s %s", r.Method, r.Path)
}

func (r *Rule) PathDepth() int {
	return len(r.Segments)
}
