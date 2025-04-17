package log_manager
import (
	file_manager "file_manager"
	"sync"
)
type LogManager struct {
	file_manager *file_manager.FileManager
	block *file_manager.BlockId
	page *file_manager.Page
	file_name string
	saved_id uint64
	tail_id uint64
	mutex sync.Mutex
}
func (log_manager *LogManager) AppendBlock(file_name string) (*file_manager.BlockId, error) {
	log_manager.mutex.Lock()
	defer log_manager.mutex.Unlock()
	block, err := log_manager.file_manager.Append(file_name)
	if err != nil {
		return nil, err
	}
	log_manager.file_manager.SetInt(0, log_manager.file_manager.Block_Size())
	log_manager.file_manager.Write(&block, log_manager.page)
	return log_manager.block, nil
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
	if file_manager.Size(file_name) <= 0 {
		block ,err := log_manager.AppendBlock(file_name)
		if err != nil {
			return nil, err
		}
		log_manager.block = block
	} else {
		block, err = log_manager.file_manager.NewBlockId(file_name, uint64(file_manager.Size(file_name) - 1))
		if err != nil {
			return nil, err
		}
		log_manager.file_manager.Read(block, log_manager.page)
		log_manager.block = block
	}
	
	 return log_manager, nil
}
func (log_manager *LogManager) Flush() error {
	log_manager.mutex.Lock()
	defer log_manager.mutex.Unlock()
	_, err := log_manager.file_manager.Write(log_manager.block, log_manager.page)
	if err != nil {
		return err
	}
	return nil
}
func (log_manager *LogManager) FlushById(id uint64) error {
	log_manager.mutex.Lock()
	defer log_manager.mutex.Unlock()
	if id > log_manager.saved_id {
		err := log_manager.Flush()
		if err != nil {
			return err
		}
		log_manager.saved_id = id
	}
	return nil
}
func (log_manager *LogManager) Append(log []byte) (uint64, error) {
	log_manager.mutex.Lock()
	defer log_manager.mutex.Unlock()
	if uint64(log_manager.file_manager.GetInt(0)) - uint64(8 + uint64(len(log))) < 8 {
		err := log_manager.Flush()
		if err != nil {
			return log_manager.tail_id, err
		}
		block, err := log_manager.AppendBlock(log_manager.file_name)
		if err != nil {
			return log_manager.tail_id, err
		}
		log_manager.block = block
	}
	log_manager.tail_id++
	log_manager.file_manager.SetBytes(uint64(log_manager.file_manager.GetInt(0)) - uint64(8 + uint64(len(log))), log)
	log_manager.file_manager.SetInt(0, uint64(log_manager.file_manager.GetInt(0)) - uint64(8 + uint64(len(log))))
	return log_manager.tail_id, nil
}
func (log_manager *LogManager) Iterator() (*LogIterator, error) {
	log_manager.mutex.Lock()
	defer log_manager.mutex.Unlock()
	log_manager.Flush()
	return NewLogIterator(log_manager.file_manager, log_manager.block), nil
}
