package security

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	NUMBER   = "0123456789"
	HEX      = "0123456789abcdef"
	ALPHA    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ALNUM    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ALPHABET = "abcdefghijklmnopqrstuvwxyz"

	ChunkSize  = 4
	ChunkCount = 4 // 16 symbols
)

type Manager interface {
	CreateGuid() uint64
	DecodeGuid(id uint64) string
	EncodeGuid(id string) uint64

	CreateIdentifier(chars string, chunks uint8) string
	DecodeIdentifier(id string, chunks uint8) string
	EncodeIdentifier(id string, chars string, chunks uint8) string

	Filter(id, chars string) string
}

type engine struct{}

var m Manager = &engine{}

func (e engine) CreateGuid() uint64 {
	return rand.Uint64()
}

func (e engine) DecodeGuid(id uint64) string {
	s := strconv.FormatUint(id, 16)
	return e.DecodeIdentifier(s, ChunkCount)
}

func (e engine) EncodeGuid(id string) uint64 {
	s := e.EncodeIdentifier(strings.ToLower(id), HEX, ChunkCount)
	res, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0
	}
	return res
}

func (e engine) Filter(id, chars string) string {
	exp, err := regexp.Compile("[^" + chars + "]")
	if err != nil {
		return id
	} else {
		return string(exp.ReplaceAll([]byte(id), []byte{}))
	}
}

func (e engine) CreateIdentifier(chars string, chunks uint8) string {
	return String(chars, chunks*ChunkSize)
}

// Humanize identifier
func (e engine) DecodeIdentifier(id string, chunks uint8) string {
	l := utf8.RuneCountInString(id)
	c := int(chunks * ChunkSize)
	if l < c {
		id = dupeChar('0', c-l) + id
	}

	var src int = 0
	var dst int = 0
	var res = make([]byte, c+int(chunks)-1)
	lim := int(chunks)
	for i := 0; i < lim; i++ {
		if i != 0 {
			res[dst] = '-'
			dst++
		}
		for j := 0; j < ChunkSize; j++ {
			res[dst] = id[src]
			src++
			dst++
		}
	}
	return string(res)
}

func (e engine) EncodeIdentifier(id string, chars string, chunks uint8) string {
	id = e.Filter(id, chars)
	maxLen := chunks * ChunkSize
	if len(id) > int(maxLen) {
		id = id[:maxLen]
	}
	return id
}

func dupeChar(ch byte, count int) string {
	res := make([]byte, count)
	for i := 0; i < count; i++ {
		res[i] = ch
	}
	return string(res)
}

func New() Manager {
	return m
}

func String(chars string, length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Int63()%int64(len(chars))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
