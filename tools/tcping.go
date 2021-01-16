package tools

import (
	"net"
	"strconv"
	"sync"
	"time"
)

func Dial(address string, timeout time.Duration) error {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func tcping(host string, port int) (t float32, ok bool) {
	start := time.Now()
	err := Dial(host+":"+strconv.Itoa(port), 3*time.Second)
	elapsed := time.Since(start)
	if err != nil {
		return 0, false
	} else {
		i := float32(elapsed.Nanoseconds()) / 1e6
		return i, true
	}
}

func tcpingAdapter(host string, port int, list *[]float32, do func()) {
	defer do()
	t, ok := tcping(host, port)
	if ok {
		*list = append(*list, t)
	}
}

func Tcping(host string, port int, count int) (max float32, min float32, avg float32, loss int) {
	var wg sync.WaitGroup
	times := new([]float32)
	wg.Add(count)
	tcping(host, port)
	for i := 0; i < count; i++ {
		go tcpingAdapter(host, port, times, wg.Done)
		time.Sleep(time.Millisecond * 20)
	}
	wg.Wait()
	if len(*times) == 0 {
		return -1, -1, -1, count
	}
	MIN := (*times)[0]
	MAX := (*times)[0]
	var SUM float32
	for _, x := range *times {
		if x > MAX {
			MAX = x
		}
		if x < MIN {
			MIN = x
		}
		SUM += x
	}
	return MAX, MIN, SUM / float32(len(*times)), count - len(*times)
}

func Go_Tcping(host string, port int, count int, c chan float32) {
	_, _, avg, _ := Tcping(host, port, count)
	c <- avg
}
