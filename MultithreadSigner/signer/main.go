package main

import (
	"strconv"
	"runtime"
	"sync"
	"strings"
	"sort"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	if len(jobs) < 1 {
		return
	}

	for _,j := range jobs[:len(jobs) - 1] {
		captureIn, captureJob := in, j
		out := make(chan interface{})

		go func() {
			defer close(out)
			captureJob(captureIn, out)
		}()
		in = out
	}
	fakeOut := make(chan interface{})
	jobs[len(jobs) - 1](in, fakeOut)
	close(fakeOut)
}

func SingleHash(in, out chan interface{}) {
	singleHashesWait := &sync.WaitGroup{}
	for val := range in {
		data := strconv.Itoa(val.(int))
		md5 := DataSignerMd5(data)
		singleHashesWait.Add(1)
		go func () {
			defer singleHashesWait.Done()

			hashOut, hashMD5Out := make(chan string), make(chan string)
			go func () {
				hashOut<- DataSignerCrc32(data)
			}()

			go func () {
				hashMD5Out<- DataSignerCrc32(md5)
			}()

			out<- <-hashOut + "~" + <-hashMD5Out
			close(hashOut)
			close(hashMD5Out)
		}()
	}
	singleHashesWait.Wait()
}

const (
	numHashes = 6
)

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for val := range in {
		captureVal := val
		wg.Add(1)
		go func () {
			defer wg.Done()
			res := make([]string, numHashes)
			var outs [numHashes]chan string
			for th, _ := range outs {
				captureTh := th
				outs[th] = make(chan string)
				go func() {
					outs[captureTh]<- DataSignerCrc32(strconv.Itoa(captureTh) + captureVal.(string))
				}()
			}
			runtime.Gosched()

			ready := 0
	MULTIHASHLOOP:
			for {
				for idx, hash := range outs {
					select {
						case x, ok := <-hash:
							if ok {
								res[idx] = x
								close(outs[idx])
								ready++
							}
						default:
						}
				}
				if ready == numHashes {
					break MULTIHASHLOOP
				}
			}
			out<- strings.Join(res, "")
		}()
	}
	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	var res []string
	for str := range in {
		res = append(res, str.(string))
	}
	sort.Strings(res)
	out<- strings.Join(res, "_")
}