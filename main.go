package main

import (
	"context"
	"fmt"
	"github.com/tvdburgt/go-argon2"
	"os"
	"strconv"
	"time"
)

const MemTestLen int = 4

type bmResult struct {
	Memory uint32 // Memory in megabytes
	Passes int    // 3 - 70
	Time   int64  // Time in milliseconds
}

// submitTime prints time difference between startTimer and the end of the parent function's execution.
func submitTime(startTimer time.Time, timeCh chan time.Duration) {
	timeCh <- time.Since(startTimer)
}

func hash(conf *argon2.Context, in []byte, start time.Time, timeCh chan time.Duration, resCh chan error) {
	defer close(resCh)
	defer submitTime(start, timeCh)

	_, err := argon2.HashEncoded(conf, in, []byte("SomeSecretSalt16"))
	resCh <- err
}

func getHashOrTimeout(ctx context.Context, conf *argon2.Context,
	in []byte, timeCh chan time.Duration) error {
	res := make(chan error)

	go func() {
		hash(conf, in, time.Now(), timeCh, res)
	}()

	for {
		select {
		case dst := <-res:
			return dst
		case <-ctx.Done():
			go func() { timeCh <- 0 }()
			return ctx.Err()
		}
	}
}

func generateConfig(passes uint32, memory uint32, threads uint8) argon2.Context {
	return argon2.Context{
		HashLen:     32,
		Iterations:  int(passes),
		Memory:      int(memory),
		Parallelism: int(threads),
		Mode:        argon2.ModeArgon2i,
		Version:     argon2.Version13,
	}
}

// low-level benchmarking func,
// uses the loop count for passes and the mem param for memory count.
func run(mem uint32, in []byte, threads uint8, maxTime uint32, runs uint8) bmResult {
	var totalTime time.Duration
	var timeCh = make(chan time.Duration)
	var kill bool

	var localResult bmResult

	for i := 3; i <= 70; i++ {
		for j := 0; uint8(j) < runs; j++ {
			func() {
				dummy := in
				conf := generateConfig(uint32(i), mem, threads)

				t := time.Duration(maxTime) * time.Millisecond
				ctx, cancel := context.WithTimeout(context.Background(), t)
				defer cancel()

				err := getHashOrTimeout(ctx, &conf, dummy, timeCh)
				if err != nil {
					kill = true
					return
				}
			}()

			totalTime += <-timeCh

			if kill {
				break
			}
		}
		if totalTime != 0 {
			endMem := mem / 1024
			endTIme := totalTime.Milliseconds() / int64(runs)
			fmt.Printf("%dmb\t%d\t%dms\n", endMem, i, endTIme)
			totalTime = 0

			localResult = bmResult{
				Memory: endMem,
				Passes: i,
				Time:   endTIme,
			}
		}

		if kill {
			close(timeCh)
			break
		}
	}

	return localResult
}

func memtest() [MemTestLen]uint32 {
	return [MemTestLen]uint32{64 * 1024, 128 * 1024, 256 * 1024, 512 * 1024}
}

// Top-level func for benchmarking execution times for hashing using:
//
//
// - Memory for thread pool (64mb/128mb/256mb/512mb)
//
// - Iterations or passes (3 - 70).
func benchmark(testData []byte, threads uint8, maxTime uint32, runs uint8) {
	topRuns := [MemTestLen]bmResult{}
	for i, mem := range memtest() {
		topRuns[i] = run(mem, testData, threads, maxTime, runs)
	}

	fmt.Println("\nLongest runs:")
	fmt.Println("MEMORY\tPASS\tTIME")
	for _, run := range topRuns {
		if run.Memory != 0 {
			fmt.Printf("%dmb\t%d\t%dms\n", run.Memory, run.Passes, run.Time)
		}
	}
}

func main() {

	var threads uint64
	if threadStr, found := os.LookupEnv("threads"); found {
		if threadParsed, err := strconv.ParseUint(threadStr, 10, 8); err != nil || threadParsed == 0 {
			panic("Invalid assignment of threads!")
		} else {
			threads = threadParsed
		}
	} else {
		threads = 1
	}

	var maxTime uint64
	if maxTimeStr, found := os.LookupEnv("maxtime"); found {
		if maxTimeParsed, err := strconv.ParseUint(maxTimeStr, 10, 32); err != nil || maxTimeParsed == 0 {
			panic("Invalid assignment of maximum run time!")
		} else {
			maxTime = maxTimeParsed
		}
	} else {
		maxTime = 500
	}

	var runs uint64
	if runsStr, found := os.LookupEnv("runs"); found {
		if runsParsed, err := strconv.ParseUint(runsStr, 10, 8); err != nil || runsParsed == 0 {
			panic("Invalid assignment of benchmark runs!")
		} else {
			runs = runsParsed
		}
	} else {
		runs = 1
	}

	fmt.Printf("Threads = %d; Maximum Time = %d; Number of runs = %d;\n", threads, maxTime, runs)
	fmt.Println("MEMORY\tPASSES\tTIME\t")
	benchmark([]byte("test!input"), uint8(threads), uint32(maxTime), uint8(runs))
}
