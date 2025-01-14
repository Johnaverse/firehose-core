package blockpoller

import (
	"context"
	"fmt"
	"testing"
	"time"

	pbbstream "github.com/streamingfast/bstream/pb/sf/bstream/v1"
	"github.com/streamingfast/firehose-core/rpc"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type TestBlockItem struct {
	skipped bool
	err     error
	delay   time.Duration
	block   *pbbstream.Block
}
type TestBlockClient struct {
	currentIndex int
	blockItems   []*TestBlockItem

	blockProduceCount int
	skippedCount      int
	errProduceCount   int
	name              string
}

func NewTestBlockClient(blockItems []*TestBlockItem, name string) *TestBlockClient {
	return &TestBlockClient{
		blockItems: blockItems,
		name:       name,
	}
}

func (c *TestBlockClient) GetBlock(ctx context.Context, blockNumber uint64) (*TestBlockItem, error) {
	fmt.Printf("%s: GetBlock %d\n", c.name, blockNumber)
	if c.EOB() {
		fmt.Println("TestBlockClient: EOB", blockNumber)
		return nil, TestErrCompleteDone
	}
	b := c.blockItems[c.currentIndex]
	fmt.Printf("%s: current index %d, block %d\n", c.name, c.currentIndex, b.block.Number)
	if b.block.Number != blockNumber {
		panic(fmt.Sprintf("%s expected requested block %d, got %d", c.name, b.block.Number, blockNumber))
	}

	c.currentIndex++
	if b.delay > 0 {
		fmt.Printf("%s: delaying block %d for %s\n", c.name, blockNumber, b.delay)
		select {
		case <-ctx.Done():
			fmt.Printf("%s: context done, returning %s\n", c.name, ctx.Err())
			return nil, ctx.Err()
		case <-time.After(b.delay):
		}
	}

	if b.err != nil {
		fmt.Printf("%s: producing block with error %d\n", c.name, blockNumber)
		c.errProduceCount++
		return nil, b.err
	}

	if b.skipped {
		fmt.Printf("%s: skipped producing block %d\n", c.name, blockNumber)
		c.skippedCount++
		return b, nil
	}

	fmt.Printf("%s: producing block %d\n", c.name, blockNumber)
	c.blockProduceCount++

	return b, nil
}

func (c *TestBlockClient) EOB() bool {
	return c.currentIndex >= len(c.blockItems)
}

type TestBlockFetcherWithClient struct {
}

func (t TestBlockFetcherWithClient) IsBlockAvailable(requestedSlot uint64) bool {
	return true
}

func (t TestBlockFetcherWithClient) Fetch(ctx context.Context, client *TestBlockClient, blkNum uint64) (b *pbbstream.Block, skipped bool, err error) {
	bi, err := client.GetBlock(ctx, blkNum)
	if err != nil {
		return nil, false, err
	}
	if bi.skipped {
		return nil, true, nil
	}
	return bi.block, false, nil
}

func TestPollerClient(t *testing.T) {
	clients := rpc.NewClients[*TestBlockClient](1*time.Second, rpc.NewRollingStrategyAlwaysUseFirst[*TestBlockClient](), zap.NewNop())
	var blockItems1 []*TestBlockItem
	var blockItems2 []*TestBlockItem

	//init call
	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("99a", "98a", 97)})
	//1st fetch block
	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("99a", "98a", 97)})

	//c1 will produce an error that c2 will be call and return requested block without error
	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("100a", "99a", 98), err: fmt.Errorf("test error")})
	blockItems2 = append(blockItems2, &TestBlockItem{block: blk("100a", "99a", 98)})

	//c1 and c2 will produce errors the c1 will be call and return requested block without error
	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("101a", "100a", 100), err: fmt.Errorf("test error")})
	blockItems2 = append(blockItems2, &TestBlockItem{block: blk("101a", "100a", 100), err: fmt.Errorf("test error")})
	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("101a", "100a", 100)})

	//test skip block
	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("102a", "101a", 101)})
	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("103a", "101a", 101), skipped: true})
	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("104a", "102a", 102)})

	blockItems1 = append(blockItems1, &TestBlockItem{block: blk("105a", "104a", 104), delay: 10 * time.Second})
	blockItems2 = append(blockItems2, &TestBlockItem{block: blk("105a", "104a", 104)})

	c1 := NewTestBlockClient(blockItems1, "c1")
	c2 := NewTestBlockClient(blockItems2, "c2")

	clients.Add(c1)
	clients.Add(c2)

	fetcher := TestBlockFetcherWithClient{}
	handler := &TestNoopBlockFinalizer{}
	poller := New(fetcher, handler, clients)

	stopBlock := uint64(106)
	err := poller.Run(99, &stopBlock, 1)

	require.NoError(t, err)

	require.Equal(t, 5, c1.blockProduceCount)
	require.Equal(t, 1, c1.skippedCount)
	require.Equal(t, 2, c1.errProduceCount)

	require.Equal(t, 2, c2.blockProduceCount)
	require.Equal(t, 0, c2.skippedCount)
	require.Equal(t, 1, c2.errProduceCount)

}
