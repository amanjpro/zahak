package book

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"

	. "github.com/amanjpro/zahak/engine"
)

type BookEntry struct {
	move   uint16
	weight uint16
}

type PolyGlot struct {
	items  map[uint64]BookEntry
	loaded bool
}

var EmptyBook = PolyGlot{nil, false}
var Book = EmptyBook

func InitBook(path string) {
	file, err := os.Open(path)

	defer file.Close()

	if err != nil {
		panic(err)
	}

	var book PolyGlot
	stat, e := file.Stat()
	if e != nil {
		book = PolyGlot{make(map[uint64]BookEntry, 100_000), false} // arbitrary initial size
	} else {
		book = PolyGlot{make(map[uint64]BookEntry, stat.Size()/16), false}
	}

	reader := bufio.NewReader(file)
	buf := make([]byte, 16)

	for {
		n, err := reader.Read(buf)

		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}

		if n == 16 {
			key := binary.LittleEndian.Uint64(buf[0:8])
			move := binary.LittleEndian.Uint16(buf[8:10])
			weight := binary.LittleEndian.Uint16(buf[10:12])
			item, found := book.items[key]
			if found {
				if item.weight < weight {
					book.items[key] = BookEntry{move, weight} // we only store the best move
				}
			} else {
				book.items[key] = BookEntry{move, weight}
			}
		} else {
			break //malformed file
		}
	}

	book.loaded = true

	Book = book
}

func GetBookMove(position *Position) Move {
	if Book.loaded {
		return EmptyMove
	}
	hash := PolyHash(position)
	entry, ok := Book.items[hash]
	if !ok {
		return EmptyMove
	}
	return ToMove(position, entry.move)
}

func ResetBook() {
	Book = EmptyBook
}
