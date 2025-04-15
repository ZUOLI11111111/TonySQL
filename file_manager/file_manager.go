package file_manager

import (
	"os"
	"path/filepath"
	"sync"
	"strings"
)

type FileManager struct {
	database_directory string
	block_size uint64
	is_new bool
	open_files map[string]*os.File
	mutex sync.Mutex
}

func NewFileManager (database_directory string, block_size uint64) (*FileManager, error) {
	file_manager := FileManager {
		database_directory: database_directory,
		block_size: block_size,
		is_new: false,
		open_files: make(map[string]*os.File),
	}
	if _, err := os.Stat(database_directory); os.IsNotExist(err) {
		file_manager.is_new = true
		err := os.Mkdir(database_directory, 0755) 
		if err != nil {
			return nil ,err
		}
	} else {
		err := filepath.Walk(database_directory, func(path string, infomation os.FileInfo, err error) error {
			if infomation.Mode().IsRegular() {
				if strings.HasPrefix(infomation.Name(), "temp") {
					os.Remove(filepath.Join(path, infomation.Name()))
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return &file_manager, nil
}

func (file_manager *FileManager) getFile(file_name string) (*os.File , error) {
	file, err := os.OpenFile(filepath.Join(file_manager.database_directory, file_name), os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	file_manager.open_files[filepath.Join(file_manager.database_directory, file_name)] = file
	return file, nil
}

func (file_manager *FileManager) Read(blockId *BlockId, page *Page) (uint64, error) {
	file_manager.mutex.Lock()
	defer file_manager.mutex.Unlock()
	file, err := file_manager.getFile(blockId.FileName())
	if err != nil {
		return 0, err
	}
	defer file.Close()
	cnt_4_byters, err := file.ReadAt(page.contents(), int64(blockId.Number()) * int64(file_manager.Block_Size()))
	if err != nil {
		return 0, err
	}
	return uint64(cnt_4_byters), nil
}

func (file_manager *FileManager) Write(blockId *BlockId, page *Page) (uint64, error) {
	file_manager.mutex.Lock()
	defer file_manager.mutex.Unlock()
	file, err := file_manager.getFile(blockId.FileName())
	if err != nil {
		return 0, err
	}
	defer file.Close()
	cnt_4_byters, err := file.WriteAt(page.contents(), int64(blockId.Number()) * int64(file_manager.Block_Size()))
	if err != nil {
		return 0, err
	}
	return uint64(cnt_4_byters), nil
}

func (file_manager *FileManager) size(file_name string) (uint64, error) {
	file, err := file_manager.getFile(file_name)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	file_content, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return uint64(file_content.Size()) / file_manager.Block_Size(), nil
}

func (file_manager *FileManager) Append(file_name string) (BlockId, error) {
	blockid, err := file_manager.size(file_name)
	if err != nil {
		return BlockId{}, err
	}
	block := NewBlockId(file_name, blockid)
	file, err := file_manager.getFile(file_name)
	if err != nil {
		return BlockId{}, err
	}
	defer file.Close()
	_, err = file.WriteAt(make([]byte, file_manager.Block_Size()), int64(blockid * file_manager.Block_Size()))
	if err != nil {
		return BlockId{}, err
	}
	return *block, nil
}

func (file_manager *FileManager) Is_New() bool {
	return file_manager.is_new
}

func (file_manager *FileManager) Block_Size() uint64 {
	return file_manager.block_size
}

	
