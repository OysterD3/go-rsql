package main

type Keyword string

const (
	SortKey           Keyword = "sort"
	FilterKey         Keyword = "filter"
	NextCursorKey     Keyword = "next"
	PreviousCursorKey Keyword = "prev"
	LimitKey          Keyword = "limit"
	OffsetKey         Keyword = "offset"
)

var ReservedCharacterMap = map[rune]bool{
	'"':  true,
	'\'': true,
	'(':  true,
	';':  true,
	',':  true,
	'=':  true,
	'!':  true,
	'>':  true,
	'<':  true,
	' ':  true,
	'\t': true,
	'\n': true,
	'\r': true,
}
