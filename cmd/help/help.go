package help

import (
	_ "embed"
)

//go:embed help.txt
var Help string

//go:embed setting.txt
var Setting string

//go:embed node.txt
var Node string

//go:embed sub.txt
var Sub string

//go:embed filter.txt
var Filter string

//go:embed recycle.txt
var Recycle string

//go:embed routing.txt
var Routing string

//go:embed alias.txt
var Alias string
