package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type NoEqd53Msg struct {
	num    uint8
	space  uint8
	writer io.Writer
}

var (
	ErrInvalidRequest = errors.New("invalid request")
)

const (
	workerIdBits = uint64(4)
	maxWorkerId  = int64(-1) ^ (int64(-1) << workerIdBits)
	sequenceBits = uint64(10)
	// This gives us over 17 years for Javascript to support integers > 53 bits
	timestampBits      = uint64(39)
	workerIdShift      = sequenceBits
	timestampLeftShift = sequenceBits + workerIdBits
	sequenceMask       = int64(-1) ^ (int64(-1) << sequenceBits)

	// Mon, 17 Jun 2013 11:48:00.000 GMT
	shepoch     = int64(1371487680000)
	idSpacesNum = int(256)
)

// Flags
var (
	wid   = flag.Int64("w", 0, "worker id")
	laddr = flag.String("l", "0.0.0.0:4444", "the address to listen on")
	lts   = flag.Int64("t", -1, "the last timestamp in milliseconds")
)

var (
	idSpacesSeq  [idSpacesNum]int64
	idSpacesMutex [idSpacesNum]sync.Mutex
)

func main() {
	parseFlags()
	acceptAndServe(mustListen())
}

func parseFlags() {
	flag.Parse()
	if *wid < 0 || *wid > maxWorkerId {
		log.Fatalf("worker id must be between 0 and %d", maxWorkerId)
	}
}

func mustListen() net.Listener {
	l, err := net.Listen("tcp", *laddr)
	if err != nil {
		log.Fatal(err)
	}
	return l
}

func acceptAndServe(l net.Listener) {
	for {
		cn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}

		go func() {
			err := serve(cn, cn)
			if err != nil {
				log.Println(err)
			}
			cn.Close()
		}()
	}
}

func serve(r io.Reader, w io.Writer) error {
	c := make([]byte, 2)
	for {
		// Wait for 2 byte request (num_ids, id_space)
		_, err := io.ReadFull(r, c)
		if err != nil && err == io.EOF {
			return nil
		}

		if c[0] == 0 {
			return ErrInvalidRequest
		}

		n := &NoEqd53Msg{c[0], c[1], w}
		n.Process()
	}

	panic("not reached")
}

func (n *NoEqd53Msg) Process() error {
	for i := uint8(0); i < n.num; i++ {
		id, err := nextId(&idSpacesSeq[n.space], &idSpacesMutex[n.space])
		if err != nil {
			return err
		}

		err = binary.Write(n.writer, binary.BigEndian, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func nextId(seq *int64, lock *sync.Mutex) (int64, error) {
	ts := milliseconds()

	if ts < *lts {
		return 0, fmt.Errorf("time is moving backwards, waiting until %d\n", *lts)
	}

	lock.Lock()
	defer lock.Unlock()
	if *lts == ts {
		*seq = (*seq + 1) & sequenceMask
		if *seq == 0 {
			for ts <= *lts {
				ts = milliseconds()
			}
		}
	} else {
		*seq = 0
	}

	*lts = ts
	ts -= shepoch
	if ts & (int64(-1) & (int64(-1) << timestampBits)) != 0 {
		panic("max timestamp value reached!")
	}

	id := (ts << timestampLeftShift) |
		(*wid << workerIdShift) |
		*seq

	return id, nil
}
