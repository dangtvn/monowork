package monowork

import (
	"os"
	"sync"
)

type Reader struct {
	InitialFrame int
	UnitFrame    int

	Index int
	File  *os.File

	Store          *sync.Map
	BufferStoreKey string
	InfoStoreKey   string

	Lock sync.RWMutex
}

type ReaderStoreData struct {
	InitialBuffer []byte
	UnitBuffer    []byte
	Timeout       int
	Order         int
}

type Song struct {
	Artist string
	Title  string
}

func NewReader() *Reader {
	return &Reader{
		InitialFrame: 0,
		UnitFrame:    0,
	}
}

func (r *Reader) GetSongInfo() *Song {
	return &Song{
		Artist: "Unknown",
		Title:  "Unknown",
	}
}
