package whoops

import (
	"fmt"
	"strings"
)

type Group []error

func (g *Group) Add(err error) {
	// nil errors are ignored
	if err == nil {
		return
	}
	*g = append(*g, err)
}

func (g Group) Err() bool {
	return len(g) > 0
}

func (g Group) Error() string {
	if len(g) == 0 {
		return "no errors in the group"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("there are %d errors in the group:\n", len(g)))

	for i := range g {
		sb.WriteString(fmt.Sprintf("error %2d: %s\n", i+1, g[i]))
	}

	return sb.String()
}

type iteratorGroupUnwrapper struct {
	g Group
	i int
}

func (i iteratorGroupUnwrapper) Is(target error) bool {
	return Is(i.g[i.i], target)
}

func (i iteratorGroupUnwrapper) As(target any) bool {
	return As(i.g[i.i], target)
}

func (i iteratorGroupUnwrapper) Error() string {
	return i.g[i.i].Error()
}

func (i iteratorGroupUnwrapper) Unwrap() error {
	if i.i+1 >= len(i.g) {
		return nil
	}
	return iteratorGroupUnwrapper{
		i: i.i + 1,
		g: i.g,
	}
}

func (g Group) Unwrap() error {
	if len(g) == 0 {
		return nil
	}
	return iteratorGroupUnwrapper{
		g: g[:],
	}
}
