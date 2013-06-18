package main

import (
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
	notify chan bool
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
	idSpacesNum = int(255)
)

// Flags
var (
	wid   = flag.Int64("w", 0, "worker id")
	laddr = flag.String("l", "0.0.0.0:4444", "the address to listen on")
	lts   = flag.Int64("t", -1, "the last timestamp in milliseconds")
)

var (
	idSpacesSeq  [idSpacesNum]int64
	idSpacesChan [idSpacesNum]chan *NoEqd53Msg
)

func init() {
	for i := range idSpacesChan {
		idSpacesChan[i] = make(chan *NoEqd53Msg)
	}
}

func main() {
	parseFlags()
	acceptAndServe(mustListen())
}

func startHandlers() {
	for i := range idSpacesChan {
		go func(c chan *NoEqd53Msg) {
			for {
				select {
				case req := <-c:
					processRequest(req)
				}
			}
		}(idSpacesChan[i])
	}
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
			err := serve(cn, cn, nil)
			if err != io.EOF {
				log.Println(err)
			}
			cn.Close()
		}()
	}
}

var once sync.Once

func serve(r io.Reader, w io.Writer, notify chan bool) error {
	once.Do(startHandlers)

	c := make([]byte, 2)
	for {
		// Wait for 2 byte request (num_ids, id_space)
		_, err := io.ReadFull(r, c)
		if err != nil {
			return err
		}

		n := uint(c[0])
		if n == 0 {
			return ErrInvalidRequest
		}

		idSpacesChan[c[1]] <- &NoEqd53Msg{c[0], c[1], w, notify}
	}

	panic("not reached")
}

func processRequest(req *NoEqd53Msg) error {
	defer func() {
		if req.notify != nil {
			req.notify <- true
		}
	}()

	b := make([]byte, uint32(req.num)*8)
	for i := uint8(0); i < req.num; i++ {
		id, err := nextId(&idSpacesSeq[req.space])
		if err != nil {
			return err
		}

		off := i * 8
		b[off+0] = byte(id >> 56)
		b[off+1] = byte(id >> 48)
		b[off+2] = byte(id >> 40)
		b[off+3] = byte(id >> 32)
		b[off+4] = byte(id >> 24)
		b[off+5] = byte(id >> 16)
		b[off+6] = byte(id >> 8)
		b[off+7] = byte(id)
	}

	_, err := req.writer.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func nextId(seq *int64) (int64, error) {
	ts := milliseconds()

	if ts < *lts {
		return 0, fmt.Errorf("time is moving backwards, waiting until %d\n", *lts)
	}

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

	id := ((ts - shepoch) << timestampLeftShift) |
		(*wid << workerIdShift) |
		*seq

	return id, nil
}
