package cosmos

import (
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type BlockService struct {
	Latest     *Block
	Blocks     map[int]*Block
	m          sync.RWMutex
	httpClient *HTTPClient
}

func NewBlockService(httpClient *HTTPClient) (*BlockService, error) {
	s := &BlockService{
		Blocks:     make(map[int]*Block),
		httpClient: httpClient,
	}

	block, err := s.httpClient.GetBlock(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.WriteBlock(block, true)

	return s, nil
}

func (s *BlockService) WriteBlock(block *Block, latest bool) {
	s.m.Lock()
	if latest {
		s.Latest = block
	}
	s.Blocks[block.Height] = block
	s.m.Unlock()
}

func (s *BlockService) GetBlock(height int) (*Block, error) {
	block, ok := s.Blocks[height]

	if ok {
		return block, nil
	}

	block, err := s.httpClient.GetBlock(&height)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.WriteBlock(block, false)

	return block, nil
}

func (c *HTTPClient) GetBlock(height *int) (*Block, error) {
	var res *BlockResponse
	var resErr *struct {
		RPCErrorResponse
		Message string `json:"message"`
	}
	req := c.tendermint.R().SetResult(&res).SetError(&resErr)

	if height != nil {
		req.SetQueryParams(map[string]string{"height": strconv.Itoa(*height)})
	}

	_, err := req.Get("/block")

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get block: %d", height)
	}

	if resErr != nil {
		if resErr.Message != "" {
			return nil, errors.Errorf("failed to get block: %d: %s", height, resErr.Message)
		}

		return nil, errors.Errorf("failed to get block: %d: %s", height, resErr.Error.Data)
	}

	if res == nil {
		return nil, errors.Errorf("res is nil for height: %d", height)
	}
	if res.Result == nil {
		return nil, errors.Errorf("res.Result is nil for height: %d", height)
	}
	if res.Result.Block == nil {
		return nil, errors.Errorf("res.Result.Block is nil for height: %d", height)
	}

	timestamp, err := time.Parse(time.RFC3339, res.Result.Block.Header.Time)
	if err != nil {
		logger.Errorf("failed to parse timestamp: %s", res.Result.Block.Header.Time)
		timestamp = time.Now()
	}

	h, err := strconv.Atoi(res.Result.Block.Header.Height)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert block height: %s", res.Result.Block.Header.Height)
	}

	b := &Block{
		Height:    h,
		Hash:      res.Result.BlockID.Hash,
		Timestamp: int(timestamp.Unix()),
	}

	return b, nil
}
