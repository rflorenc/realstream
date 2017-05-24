package main

import "strings"

/*

This simple parser provides a NewPath function that parses the specified path string
and returns a new instance of the Path type. Leading and trailing slashes are trimmed
(using strings.Trim ) and the remaining path is split (using strings.Split ) by the
PathSeparator constant that is just a forward slash. If there is more than one segment
( len(s) > 1 ), the last one is considered to be the ID. We re-slice the slice of strings
to select the last item for the ID using s[len(s)-1] , and the rest of the items for the
remainder of the path using s[:len(s)-1] . On the same lines, we also re-join the path
segments with the PathSeparator constant to form a single string containing the path
without the ID.
This supports any collection/id pair, which is all we need for our API.

Gorilla "Mux" package takes care of this but for our simple use case
we don't need to add any extra dependencies.

*/

const PathSeparator = "/"

type Path struct {
	Path string
	ID   string
}

func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, PathSeparator)
	s := strings.Split(p, PathSeparator)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparator)
	}
	return &Path{Path: p, ID: id}
}

func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
