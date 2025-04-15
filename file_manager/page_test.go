package file_manager

import (
	"testing"
	"github.com/stretchr/testify/require"
	"fmt"
)

func TestSetAndGetInt(t *testing.T) {
	//require.Equal(t, 1, 2)
	page := NewPageBySize(256)
	val := uint64(1234)
	offset := uint64(23)
	page.SetInt(offset, val)
	val_get := page.GetInt(offset)
	require.Equal(t, val, val_get)
}

func TestSetAndGetByteArray(t *testing.T) {

	page := NewPageBySize(256)
	bs := []byte{1, 2, 3, 4, 5, 6}
	offset := uint64(111)
	page.SetBytes(offset, bs)
	bs_got := page.GetBytes(offset)

	require.Equal(t, bs, bs_got)
	fmt.Println(bs_got)
}

func TestSetAndGetString(t *testing.T) {
	page := NewPageBySize(256)
	str := "hello world"
	offset := uint64(111)
	page.SetString(offset, str)
	str_got := page.GetString(offset)
	require.Equal(t, str, str_got)
	fmt.Println(str_got)
}

func TestMaxLengthForString(t *testing.T) {
	page := NewPageBySize(256)
	str := "hello world"
	max_length := page.MaxLengthForString(str)
	require.Equal(t, 8 + uint64(len(str)), max_length)
	fmt.Println(max_length)
}

func TestSetStringWithMaxLength(t *testing.T) {
	page := NewPageBySize(256)
	str := "hello world"
	offset := uint64(111)
	max_length := page.MaxLengthForString(str)
	page.SetString(offset, str)
	str_got := page.GetString(offset)
	require.Equal(t, str, str_got)
	fmt.Println(str_got)
	fmt.Println(max_length)
}
func TestContents(t *testing.T) {
	page := NewPageBySize(256)
	page.SetString(111, "hello world")
	content := page.contents()
	fmt.Println(content)
}
