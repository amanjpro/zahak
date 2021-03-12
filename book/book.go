package book

import (
	"bufio"
	"encoding/binary"
	"io"
	"math/rand"
	"os"
	"time"

	. "github.com/amanjpro/zahak/engine"
)

type PolyGlot struct {
	items  map[uint64][]uint16
	loaded bool
}

var EmptyBook = PolyGlot{nil, false}
var Book = EmptyBook

func InitBook(path string) {
	rand.Seed(time.Now().Unix())
	file, err := os.Open(path)

	defer file.Close()

	if err != nil {
		panic(err)
	}

	var book PolyGlot
	stat, e := file.Stat()
	if e != nil {
		book = PolyGlot{make(map[uint64][]uint16, 100_000), false} // arbitrary initial size
	} else {
		book = PolyGlot{make(map[uint64][]uint16, stat.Size()/16), false}
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
			key := binary.BigEndian.Uint64(buf[0:8])
			move := binary.BigEndian.Uint16(buf[8:10])
			item, found := book.items[key]
			if found {
				book.items[key] = append(item, move)
			} else {
				item = make([]uint16, 0, 30)
				item = append(item, move)
				book.items[key] = item
			}
		} else {
			break //malformed file
		}
	}

	book.loaded = true

	Book = book
}

func GetBookMove(position *Position) Move {
	if !Book.loaded {
		return EmptyMove
	}
	hash := PolyHash(position)
	items, ok := Book.items[hash]
	if !ok {
		return EmptyMove
	}
	index := rand.Intn(len(items))
	return ToMove(position, items[index])
}

func ResetBook() {
	Book = EmptyBook
}

func IsBoookLoaded() bool {
	return Book.loaded
}
