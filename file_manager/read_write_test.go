package file_manager

import (
	"testing"
	"github.com/stretchr/testify/require"
	"os"
)

func TestReadWrite(t *testing.T) {
	// 清理测试目录
	os.RemoveAll("read_write_test")
	
	// 创建新的 FileManager
	fm, err := NewFileManager("read_write_test", 400)
	require.NoError(t, err)
	
	// 创建 BlockId 和 Page
	blk := NewBlockId("testfile", 2)
	p1 := NewPageBySize(fm.Block_Size())
	
	// 在 Page 中设置值
	val := uint64(345)
	p1.SetInt(100, val)
	
	// 写入数据
	_, err = fm.Write(blk, p1)
	require.NoError(t, err)
	
	// 创建新的 Page 并读取数据
	p2 := NewPageBySize(fm.Block_Size())
	_, err = fm.Read(blk, p2)
	require.NoError(t, err)
	
	// 验证读取的值和写入的值是否相同
	require.Equal(t, val, p2.GetInt(100))
} 