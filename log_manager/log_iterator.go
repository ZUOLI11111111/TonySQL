package log_manager

import (
	fm "file_manager"
)

type LogIterator struct {
	file_manager *fm.FileManager
	block    *fm.BlockId
	page         *fm.Page
	bound        uint64
	current      uint64
}

func NewLogIterator(file_manager *fm.FileManager, block_id *fm.BlockId) *LogIterator {
	log_iterator := &LogIterator{
		file_manager: file_manager,
		block:     block_id,
	}
	log_iterator.page = file_manager.NewPageByBlockId(block_id.Size())
	err := log_iterator.MoveToBlock(block_id)
	return log_iterator
}

func (log_iterator *LogIterator) MoveToBlock(block *fm.Block_Size) error {
	_, err := log_iterator.log_manager.Read(block, log_iterator.page)
	if err != nil {
		return err
	}
	log_iterator.bound = log_iterator.page.GetInt(0)
	log_iterator.current = log_iterator.bound
	return nil
}

func (log_iterator *LogIterator) HasNext() bool {
	return log_iterator.current < log_iterator.bound || log_iterator.page.GetInt(log_iterator.current) >= 0
}

func (log_iterator *LogIterator) Next() ([]byte, error) {
	if !log_iterator.HasNext() {
		return nil, errors.New("no more logs")
	}
	if log_iterator.current >= log_iterator.bound {
		log_iterator.block = log_iterator.block.NewBlockId(log_iterator.block.FileName(), log_iterator.block.Number() - 1)
		err := log_iterator.MoveToBlock(log_iterator.block)
		if err != nil {
			return nil, err
		}
	}
	tmp := log_iterator.current
	log_iterator.current += 8 + uint64(log_iterator.page.GetInt(log_iterator.current))
	return log_iterator.page.GetBytes(tmp), nil
	

}


