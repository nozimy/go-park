package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
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
	wg := &sync.WaitGroup{}

	for _, doJob := range jobs {
		wg.Add(1)
		mu := &sync.Mutex{}

		mu.Lock()
		go func(inInner, outInner chan interface{}, doJobInner job) {
			defer mu.Unlock()
			defer wg.Done()
			defer close(outInner)

			doJobInner(inInner, outInner)
		}(in, out, doJob)

		in = out
		out = make(chan interface{})
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	mu := &sync.Mutex{}
	for i := range in {
		data := strconv.Itoa(i.(int))
		md5 := make(chan string)
		crc32Md5 := make(chan string)
		crc32 := make(chan string)

		go func() {
			mu.Lock()
			md5 <- DataSignerMd5(data)
			mu.Unlock()
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
	}
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for i := range in {
		data := i.(single)
		wg.Add(1)
		go doMultiHash(data, out, wg)
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var data []string

	for i := range in {
		data = append(data, i.(string))
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	out <- strings.Join(data, "_")
}

func doDataSignerCrc32(index int, data string, resultMap map[int]string, mu *sync.Mutex, waiter *sync.WaitGroup) {
	defer waiter.Done()
	th := strconv.Itoa(index)
	resString := DataSignerCrc32(th + data)
	mu.Lock()
	resultMap[index] = resString
	mu.Unlock()
}

func doMultiHash(data single, out chan interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()

	dataStr := data.get()

	wg := &sync.WaitGroup{}
	resultMap := make(map[int]string)
	mu := &sync.Mutex{}
	var result string

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go doDataSignerCrc32(i, dataStr, resultMap, mu, wg)
	}

	wg.Wait()

	for i := 0; i < len(resultMap); i++ {
		result += resultMap[i]
	}

	out <- result
}
