package probe

import (
	"fmt"
	"strings"
)

// ResultFilter describes an interface for checking
// whether a Result meets given implementation-specific
// criteria.
type ResultFilter interface {
	Check(*Result) bool
}

// FilterResponseCode filters on a response Code equaling the int
type FilterResponseCode int

func (f FilterResponseCode) String() string {
	return fmt.Sprintf("<ResponseCode=%d>", f)
}

// Check filters a response Code when it is equal
func (f FilterResponseCode) Check(res *Result) bool {
	return res.Code == int(f)
}

// FilterResponseContains filters on the response Body having the
// string in its contents
type FilterResponseContains string

func (f FilterResponseContains) String() string {
	return fmt.Sprintf("<ResponseBody⊃%s>", string(f))
}

// Check filters when the Body contains the string
func (f FilterResponseContains) Check(res *Result) bool {
	return strings.Contains(res.Body, string(f))
}

// FilterGroupAll is true when all the member ResultFilter
// checks are true
type FilterGroupAll struct {
	Members []ResultFilter
}

func (f FilterGroupAll) String() string {
	return fmt.Sprintf("&(%s)", f.Members)
}

// Check is true when all the member ResultFilter
// checks are true
func (fg FilterGroupAll) Check(res *Result) bool {
	// Check that **ALL** member filters are true
	for _, x := range fg.Members {
		if !x.Check(res) {
			return false
		}
	}
	return true
}

// FilterGroupAny is true when any of the member ResultFilter
// checks are true
type FilterGroupAny struct {
	Members []ResultFilter
}

func (f FilterGroupAny) String() string {
	return fmt.Sprintf("|(%s)", f.Members)
}

// Check is true when any of the member ResultFilter
// checks are true
func (fg FilterGroupAny) Check(res *Result) bool {
	// Check that **ANY** member filters are true
	for _, x := range fg.Members {
		if x.Check(res) {
			return true
		}
	}
	return false
}

// FilterGroupNot only checks one member is false
//
// To build multiple not conditions, have a child
// FilterGroupAll or FilterGroupAny as a member
type FilterGroupNot struct {
	Member ResultFilter
}

func (f FilterGroupNot) String() string {
	return fmt.Sprintf("!(%s)", f.Member)
}

// FilterGroupNot is true when the member checks false
func (fg FilterGroupNot) Check(res *Result) bool {
	return !fg.Member.Check(res)
}
