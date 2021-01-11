package rule

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

type _SegmentRule struct {
	// Segment value
	value string
	// Segment pattern
	pattern *regexp.Regexp
	// Variable name
	varName string
}

type Rule struct {
	// Path depth
	depth int
	// Path segment rules
	srs []*_SegmentRule
}

// Match tests does the given path match the rule, returns the matching score.
// The higher score means more exactly matching occurs, and the negative score means mismatching.
func (r *Rule) Match(path string, vars map[string]string) (score int) {
	// Split URL path into segments
	runes, segments := []rune(path), make([]string, 0)
	last, count := -1, len(runes)
	for i := 0; i <= count; i++ {
		if i == count || runes[i] == '/' {
			if i > 0 {
				segment := string(runes[last+1 : i])
				if decodedSegment, err := url.PathUnescape(segment); err == nil {
					segments = append(segments, decodedSegment)
				} else {
					segments = append(segments, segment)
				}
			}
			last = i
		}
	}
	if len(segments) != r.depth {
		return -1
	}
	for i := 0; i < r.depth; i++ {
		if sr := r.srs[i]; sr.pattern == nil {
			if sr.value == segments[i] {
				// Exactly match, increase score
				score += 1
			} else {
				// Break on mismatching
				return -1
			}
		} else {
			if m := sr.pattern.FindStringSubmatch(segments[i]); len(m) > 0 {
				// Add path variable to map, do not increase the score
				vars[sr.varName] = m[1]
			} else {
				// Break on mismatching
				return -1
			}
		}
	}
	return
}

func (r *Rule) Path() (path string, exactly bool) {
	path, exactly = "/", true
	if r.depth > 0 {
		sb := &strings.Builder{}
		for _, sr := range r.srs {
			if sr.pattern != nil {
				exactly = false
			}
			sb.WriteString("/" + sr.value)
		}
		path = sb.String()
	}
	return
}

/*
Parse parses a method name into route rule.

The method name should be in camel-case, consists of several words.
The first word will be treated as request method name, such as "Get", "Post", etc.
The following words will be treated as path segments, e.g. "UserData" will map to "/user/data".
And there are some special keywords for special usage, please check the reference.
*/
func Parse(name string) (method string, rule *Rule, err error) {
	// Split method name into words
	runes, words := []rune(name), make([]string, 0)
	start, count := 0, len(runes)
	for i := 0; i <= count; i++ {
		if i == count || unicode.IsUpper(runes[i]) {
			if i == start {
				continue
			}
			words = append(words, strings.ToLower(string(runes[start:i])))
			start = i
		}
	}
	method = words[0]
	if len(words) > 1 {
		rule, err = (&_Parser{words: words[1:]}).Parse()
	} else {
		rule = &Rule{
			depth: 0,
		}
	}
	return
}
