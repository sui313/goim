package comet

import (
	"github.com/Terry-Mao/goim/internal/comet/conf"
	"github.com/Terry-Mao/goim/pkg/bytes"
	"github.com/Terry-Mao/goim/pkg/time"
)

// RoundOptions round options.
type RoundOptions struct {
	Timer        int
	TimerSize    int
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

// Round userd for connection round-robin get a reader/writer/timer for split big lock.
type Round struct {
	readers []bytes.Pool
	writers []bytes.Pool
	timers  []time.Timer
	options RoundOptions
}

// NewRound new a round struct.
func NewRound(c *conf.Config) (r *Round) {
	var i int
	r = &Round{
		options: RoundOptions{
			Reader:       c.TCP.Reader, //32
			ReadBuf:      c.TCP.ReadBuf, //1024
			ReadBufSize:  c.TCP.ReadBufSize,//8192
			Writer:       c.TCP.Writer, //32
			WriteBuf:     c.TCP.WriteBuf,//1024
			WriteBufSize: c.TCP.WriteBufSize,//8192
			Timer:        c.Protocol.Timer,//32
			TimerSize:    c.Protocol.TimerSize,//2048
		}}
	//Ĭ������:����32�����������أ�ÿ����������С1024*8192K=8M��ÿ���������ָ��1024�飬ÿ��8K��ÿ��֮���õ���������������
	//ÿ�λ�ȡ���п飬����Żس���
	// reader
	r.readers = make([]bytes.Pool, r.options.Reader) //������������
	for i = 0; i < r.options.Reader; i++ {
		r.readers[i].Init(r.options.ReadBuf, r.options.ReadBufSize)
	}
	// writer
	r.writers = make([]bytes.Pool, r.options.Writer)//����д������
	for i = 0; i < r.options.Writer; i++ {
		r.writers[i].Init(r.options.WriteBuf, r.options.WriteBufSize)
	}
	// timer
	r.timers = make([]time.Timer, r.options.Timer)
	for i = 0; i < r.options.Timer; i++ {
		r.timers[i].Init(r.options.TimerSize)
	}
	return
}

// Timer get a timer.
func (r *Round) Timer(rn int) *time.Timer {
	return &(r.timers[rn%r.options.Timer])
}

// Reader get a reader memory buffer.
func (r *Round) Reader(rn int) *bytes.Pool {
	return &(r.readers[rn%r.options.Reader])
}

// Writer get a writer memory buffer pool.
func (r *Round) Writer(rn int) *bytes.Pool {
	return &(r.writers[rn%r.options.Writer])
}
