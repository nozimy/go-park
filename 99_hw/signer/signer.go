package main

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type single struct {
	crc32Md5 chan string
	crc32    chan string
}

func (s *single) get() string {
	return <-s.crc32 + "~" + <-s.crc32Md5
}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	out := make(chan interface{})
	result := make(chan int)

	for i, doJob := range jobs {
		in = out
		out = make(chan interface{})

		go func(in2, out2 chan interface{}, i2 int, result2 chan<- int, doJob2 job, lastIndex int) {
			defer close(out2)

			doJob2(in2, out2)

			log.Println("done   doJob", i2)

			if i2 == lastIndex {
				result2 <- i2
			}
		}(in, out, i, result, doJob, len(jobs)-1)
	}

	<-result
}

var SingleHash = func(in, out chan interface{}) {
	for i := range in {
		data := strconv.Itoa(i.(int))
		md5 := make(chan string)
		crc32Md5 := make(chan string)
		crc32 := make(chan string)

		go func() {
			log.Println("DataSignerMd5")
			md5 <- DataSignerMd5(data) //перегрев на 1 сек
		}()
		go func() {
			crc32Md5 <- DataSignerCrc32(<-md5)
		}()
		go func() {
			crc32 <- DataSignerCrc32(data)
		}()

		out <- single{
			crc32Md5: crc32Md5,
			crc32:    crc32,
		}

		time.Sleep(11 * time.Millisecond)
	}
}

var MultiHash = func(in, out chan interface{}) {
	wgOuter := &sync.WaitGroup{}

	for i := range in {
		data := i.(single)

		wgOuter.Add(1)

		go func(data2 single, out2 chan interface{}, waiterOuter *sync.WaitGroup) {
			defer waiterOuter.Done()

			log.Println("MultiHash Start")
			d := <-data2.crc32 + "~" + <-data2.crc32Md5
			log.Println("MultiHash after data extraction", d)

			wg := &sync.WaitGroup{}
			resultMap := make(map[int]string)
			mu := &sync.Mutex{}
			var result string

			for i := 0; i < 6; i++ {
				wg.Add(1)
				go func(i2 int, waiter *sync.WaitGroup) {
					defer waiter.Done()
					th := strconv.Itoa(i2)

					res := DataSignerCrc32(th + d)
					mu.Lock()
					resultMap[i2] = res
					mu.Unlock()
				}(i, wg)
			}

			wg.Wait()

			for i := 0; i < len(resultMap); i++ {
				result += resultMap[i]
			}

			log.Println("MultiHash Finish")

			log.Println("MultiHash End")
			out2 <- result
		}(data, out, wgOuter)
	}

	wgOuter.Wait()
}

var CombineResults = func(in, out chan interface{}) {
	var data []string

	for i := range in {
		data = append(data, i.(string))
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	out <- strings.Join(data, "_")
}
