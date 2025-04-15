package file_manager

import (
	"crypto/sha256"
	"fmt"
)

type BlockId struct {
	file_name string
	blk_num uint64
}

func NewBlockId(file_name string, blk_num uint64) *BlockId {
	return &BlockId {
		file_name: file_name,
		blk_num: blk_num,
	}
}

func (b *BlockId) FileName() string {
	return b.file_name
}
func (b *BlockId) Number() uint64 {
	return b.blk_num
}

func (b *BlockId) Equal(B *BlockId) bool {
	return b.file_name == B.file_name && b.blk_num == B.blk_num
}

func asSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (b *BlockId) HashCode() string {
	return asSha256(*b)
}