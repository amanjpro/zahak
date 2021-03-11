package book

type BookEntry struct {
	key  uint64
	move int32
}

type PolyGlot struct {
	items  []BookEntry
	loaded bool
}

var book = PolyGlot{nil, false}

func InitBook(path string) {
	// load polyglot file here
}
