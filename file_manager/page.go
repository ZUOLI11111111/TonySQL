package file_manager

import (
	"encoding/binary"
)

type Page struct {
	buffer []byte
}

func NewPageBySize(block_size uint64) *Page {
	bytes := make([]byte, block_size)
	return &Page {
		buffer: bytes,
	}
}

func (p *Page) GetInt(offset uint64) uint64 {
	num := binary.LittleEndian.Uint64(p.buffer[offset : offset+8])
	return num
}

func uint64ToByteArray(val uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)
	return b
}

func (p *Page)SetInt(offset uint64, val uint64) {
	b := uint64ToByteArray(val)
	copy(p.buffer[offset:], b)
}

func (p *Page)SetBytes(offset uint64, bytes []byte) {
	length := uint64(len(bytes))
	length_buffer := uint64ToByteArray(length)
	copy(p.buffer[offset:], length_buffer)
	copy(p.buffer[offset + 8:], bytes)
}

func (p *Page)GetBytes(offset uint64) []byte {
	length := binary.LittleEndian.Uint64(p.buffer[offset : offset + 8])
	buffer := make([]byte, length)
	copy(buffer, p.buffer[offset + 8:])
	return buffer
}

func (p *Page)GetString(offset uint64) string {
	return string(p.GetBytes(offset))
}

func (p *Page)SetString(offset uint64, s string) {
	str_bytes := []byte(s)
	p.SetBytes(offset, str_bytes)
}

func (p *Page)MaxLengthForString(s string) uint64 {
	return uint64(8 + len([]byte(s)))
}

func (p *Page) contents() []byte {
	return p.buffer
}

