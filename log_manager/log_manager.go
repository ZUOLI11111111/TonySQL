package log_manager

import (
	file_manager "file_manager"
	"sync"
)

const (
	UINT64_LEN = 8
)

type LogManager struct {
	file_manager *file_manager.FileManager
	block *file_manager.BlockId
	page *file_manager.Page
	file_name string
	block_size uint64
	saved_id uint64
	tail_id uint64
	mutex sync.Mutex
}

func NewLogManager(file_manager *file_manager.FileManager, file_name string) (*LogManager, error) {
	log_manager := &LogManager{
		file_manager: file_manager,
		page: file_manager.NewPageBySize(file_manager.Block_Size()),
		saved_id: 0,
		file_name: file_name,
		block_size: file_manager.Block_Size(),
		tail_id: 0,
	}
	
	
	if err != nil {
		return nil, err
	}
	 return log_manager, nil
}
