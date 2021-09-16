package example

import "github.com/zhangsifeng92/geos/libraries/container"

var IntComparator = func(a, b interface{}) int { return container.IntComparator(a.(int), b.(int)) }
var StringComparator = func(a, b interface{}) int { return container.StringComparator(a.(string), b.(string)) }

//go:generate go install "github.com/zhangsifeng92/geos/libraries/container"
//go:generate go install "github.com/zhangsifeng92/geos/libraries/container/redblacktree"
//go:generate go install "github.com/zhangsifeng92/geos/libraries/container/treemap"
//go:generate gotemplate "github.com/zhangsifeng92/geos/libraries/container/treemap" IntStringMap(int,string,IntComparator,false)
//go:generate gotemplate "github.com/zhangsifeng92/geos/libraries/container/treemap" MultiIntStringMap(int,string,IntComparator,true)
//go:generate gotemplate "github.com/zhangsifeng92/geos/libraries/container/treemap" IntStringPtrMap(int,*string,IntComparator,false)
//go:generate gotemplate "github.com/zhangsifeng92/geos/libraries/container/treemap" MultiIntStringPtrMap(int,*string,IntComparator,true)
//go:generate gotemplate "github.com/zhangsifeng92/geos/libraries/container/treemap" StringIntMap(string,int,StringComparator,false)
//go:generate gotemplate "github.com/zhangsifeng92/geos/libraries/container/treemap" MultiStringIntMap(string,int,StringComparator,true)
//go:generate go build .
