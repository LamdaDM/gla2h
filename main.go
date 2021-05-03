package main

import (
	"context"
	"fmt"
	"github.com/matthewhartstonge/argon2"
	"os"
	"strconv"
	"time"
)

type HashRes struct {
	hashed []byte
	hashErr error
}

// ExecTime prints time difference between startTimer and the end of the parent function's execution.
func ExecTime(startTimer time.Time) { fmt.Println(time.Since(startTimer))}

// ProduceHash calls argon2.HashEncoded() on the config and input passed.
func ProduceHash(conf argon2.Config, stdin []byte, timer bool, start time.Time) HashRes {
	if timer { defer ExecTime(start) }

	if out, hashErr := conf.HashEncoded(stdin);
		hashErr != nil {
			return HashRes{nil, hashErr}
		} else { return HashRes{out, nil} }
}

func SafeProduceHash(
	ctx context.Context,
	conf argon2.Config,
	stdin []byte,
	timer bool,
	start time.Time,
	) HashRes {
	res := make(chan HashRes)
	go func() {
		res <- ProduceHash(conf, stdin, timer, start)
		close(res)
	}()

	for{
		select {
			case dst := <-res:
				return dst
			case <-ctx.Done():
				return HashRes{nil, ctx.Err()}
		}
	}
}

// GenerateConfig generates a new argon2.Config,
func GenerateConfig(passes uint32, memory uint32, threads uint8) argon2.Config {
	return argon2.Config{
		HashLength:  32,
		SaltLength:  16,
		TimeCost:    passes,
		MemoryCost:  memory,
		Parallelism: threads,
		Mode:        argon2.ModeArgon2i,
		Version:     argon2.Version13,
	}
}

// low-level benchmarking func,
// uses the loop count for passes and the memtest param for memory count.
// Starts the timer manually, so timing precision is lowered a bit
// due to possible time difference between the call of time.Now()
// and ProduceHash()
func run(memtest uint32, stdin []byte, timer bool, threads uint8){
	var startT time.Time
	var kill bool
	for i := 3; i <= 70; i++ {
		func() {
			dummy := stdin
			sinCtx, cancel := context.WithTimeout(context.Background(), 750*time.Millisecond)
			defer cancel()
			fmt.Printf("passes-%d mem-%d: ", i, memtest)
			startT = time.Now()
			if hashResult := SafeProduceHash(
				sinCtx,
				GenerateConfig(uint32(i), memtest, threads),
				dummy,
				timer,
				startT,
			); hashResult.hashErr != nil {
				kill = true
				return
			}
		}()
		if kill { break }
	}
}

// Top-level func for benchmarking execution times for hashing.
// Variables: memory for thread pool (64mb/128mb/256mb/512mb),
// iterations or passes (3 - 70).
func benchmark(stdin []byte, timer bool, threads uint8) {
	const memtest1 uint32 = 64*1024
	const memtest2 uint32 = 128*1024
	const memtest3 uint32 = 256*1024
	const memtest4 uint32 = 512*1024

	run(memtest1, stdin, timer, threads)
	run(memtest2, stdin, timer, threads)
	run(memtest3, stdin, timer, threads)
	run(memtest4, stdin, timer, threads)

}

func main() {
	if len(os.Args) < 7 {
		fmt.Fprint(os.Stderr, "argument count mismatch, requires 6 args. ")
		panic(main)
	}

	in := []byte(os.Args[1])
	var timer bool

	MEM_COST, memParseErr := strconv.ParseUint(os.Args[4], 10, 32)
	if memParseErr != nil {
		fmt.Fprint(os.Stderr, memParseErr)
		panic(memParseErr)
	}

	TIME_COST, timeParseErr := strconv.ParseUint(os.Args[5], 10, 32)
	if timeParseErr != nil {
		fmt.Fprint(os.Stderr, timeParseErr)
		panic(timeParseErr)
	}

	THREAD_AMOUNT, threadsParseErr := strconv.ParseUint(os.Args[6], 10, 8)
	if threadsParseErr != nil {
		fmt.Fprint(os.Stderr, threadsParseErr)
		panic(threadsParseErr)
	}

	if os.Args[2] == "y" { timer = true }
	if os.Args[3] == "y" { benchmark(in, timer, uint8(THREAD_AMOUNT))
	} else {
		conf := GenerateConfig(uint32(TIME_COST), uint32(MEM_COST) * 1024, uint8(THREAD_AMOUNT))

		res := ProduceHash(conf, in, timer, time.Now())
		if res.hashErr != nil {
			if _, printErr := fmt.Fprint(os.Stderr, res.hashErr);
				printErr != nil {
				panic(printErr)
			}
		} else {
			if _, printErr := fmt.Fprint(os.Stdout, string(res.hashed));
				printErr != nil {
				panic(printErr)
			}
		}
	}
}
