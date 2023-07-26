package counter

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const subPeriodDuration = 5
const subPeriodsCount = 12
const MetricsPeriod = subPeriodDuration * subPeriodsCount
const channelsSize = 0

type entry struct {
	periodStart int64
	sum         int64
	count       int64
}

type counter struct {
	next    *entry
	last    *entry
	entries chan *entry
	sum     int64
	count   int64
	ch      chan int64
	mx      sync.RWMutex
}

var counterMap sync.Map

var timeNowUnixFunc = func() int64 {
	return time.Now().Unix()
}

var debug = os.Getenv("DEBUG_COUNTER") != ""

func RequestProcessor(line string) string {
	debugMsg("parsing:", line)
	separatorPos := strings.Index(line, " ")

	if separatorPos == -1 {
		avg, total := getCounter(line).averageAndCount()
		debugMsg("query:", line, "results:", avg, total)
		return strconv.FormatInt(avg, 10) + " " + strconv.FormatInt(total, 10)
	}

	key := line[:separatorPos]
	value, err := strconv.ParseInt(line[separatorPos+1:], 10, 64)
	if err != nil {
		debugMsg("malformed:", line)
		return "" // don't care about malformed requests for now
	}
	debugMsg("inserting:", line)
	getCounter(key).ch <- value
	return ""
}

func getCounter(key string) *counter {
	var cnt *counter
	if existing, ok := counterMap.Load(key); !ok {
		cnt = createCounter()
		debugMsg("creating counter:", key)
		counterMap.Store(key, cnt)
	} else {
		cnt = existing.(*counter)
	}
	return cnt
}

func createCounter() *counter {
	cnt := &counter{
		next:    &entry{},
		entries: make(chan *entry, subPeriodsCount),
		ch:      make(chan int64, channelsSize),
	}
	go cnt.receiver()
	return cnt
}

func (c *counter) averageAndCount() (int64, int64) {
	c.refresh()
	if c.count == 0 {
		return 0, 0
	}
	c.mx.RLock()
	avg := c.sum / c.count
	total := c.count
	c.mx.RUnlock()
	return avg, total
}

func (c *counter) receiver() {
	for value := range c.ch {
		debugMsg("received:", value)
		c.refresh()
		c.next.sum += value
		c.next.count += 1
		c.mx.Lock()
		c.sum += value
		c.count += 1
		c.mx.Unlock()
	}
}

func (c *counter) refresh() {
	ts := timeNowUnixFunc() / subPeriodDuration
	if c.next.periodStart < ts {
		c.entries <- c.next
		c.next = &entry{periodStart: ts}
		c.flush()
	}
}

func (c *counter) flush() {
	c.dropLast()
	for c.last == nil {
		select {
		case c.last = <-c.entries:
			c.dropLast()
		default:
			return
		}
	}
}

func (c *counter) dropLast() {
	if c.last == nil || c.next.periodStart-c.last.periodStart <= subPeriodsCount {
		return
	}
	c.mx.Lock()
	c.sum -= c.last.sum
	c.count -= c.last.count
	c.mx.Unlock()
	c.last = nil
}

func debugMsg(msg ...any) {
	if debug {
		log.Println(msg...)
	}
}
