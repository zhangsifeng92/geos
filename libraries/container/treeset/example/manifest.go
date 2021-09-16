package example

import "github.com/zhangsifeng92/geos/libraries/container"

var StringComparator = func(a, b interface{}) int { return container.StringComparator(a.(string), b.(string)) }

//go:generate go install "github.com/zhangsifeng92/geos/libraries/container"
//go:generate go install "github.com/zhangsifeng92/geos/libraries/container/redblacktree"
//go:generate go install "github.com/zhangsifeng92/geos/libraries/container/treeset"
//go:generate gotemplate "github.com/zhangsifeng92/geos/libraries/container/treeset" StringSet(string,StringComparator,false)
//go:generate gotemplate "github.com/zhangsifeng92/geos/libraries/container/treeset" MultiStringSet(string,StringComparator,true)
//go:generate go build .
