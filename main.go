package main

import (
	"fmt"
	"github.com/matthewhartstonge/argon2"
	"os"
	"strconv"
	"time"
)

// ExecTime prints time difference between startTimer and the end of the parent function's execution.
func ExecTime(startTimer time.Time) { fmt.Println(time.Since(startTimer))}

// ProduceHash calls argon2.HashEncoded() on the config and input passed.
func ProduceHash(conf argon2.Config, stdin []byte, timer bool, start time.Time) ([]byte, error) {
	if timer { defer ExecTime(start) }

	if out, hashErr := conf.HashEncoded(stdin);
		hashErr != nil {
			return nil, hashErr
		} else { return out, nil }
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

// Low-level benchmarking func,
// uses the loop count for passes and the memtest param for memory count.
// Starts the timer manually, so timing precision is lowered a bit
// due to possible time difference between the call of time.Now()
// and ProduceHash()
func run(memtest uint32, stdin []byte, timer bool){
	var startT time.Time
	for i := 3; i <= 70; i++ {
		dummy := stdin
		fmt.Printf("passes-%d mem-%d: ", i, memtest)
		startT = time.Now()
		if out, hashErr := ProduceHash(
			GenerateConfig(uint32(i), memtest, 12),
			dummy,
			timer,
			startT,
		); hashErr != nil {
			if _, printErr := fmt.Fprint(os.Stderr, hashErr);
				printErr != nil { panic(printErr) }
		} else { fmt.Println(string(out)) }
	}
}

// Top-level func for benchmarking execution times for hashing.
// Variables: memory for thread pool (64mb/128mb/256mb/512mb),
// iterations or passes (3 - 70).
func benchmark(stdin []byte, timer bool) {
	const memtest1 uint32 = 64
	const memtest2 uint32 = 128
	const memtest3 uint32 = 256
	const memtest4 uint32 = 512

	run(memtest1, stdin, timer)
	run(memtest2, stdin, timer)
	run(memtest3, stdin, timer)
	run(memtest4, stdin, timer)

}

func main() {
	in := []byte(os.Args[1])
	var timer bool
	if len(os.Args) < 7 {
		fmt.Fprint(os.Stderr, "argument count mismatch, requires 6 args. ")
		panic(main)
	}
	if os.Args[2] == "timed" { timer = true }
	if os.Args[3] == "benchmark" { benchmark(in, timer)
	} else {
		MEM_COST, memParseErr := strconv.ParseUint(os.Args[4], 10, 32);
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

		conf := GenerateConfig(uint32(TIME_COST), uint32(MEM_COST), uint8(THREAD_AMOUNT))

		out, hashErr := ProduceHash(conf, in, timer, time.Now())
		if hashErr != nil {
			if _, printErr := fmt.Fprint(os.Stderr, hashErr);
				printErr != nil {
				panic(printErr)
			}
		} else {
			if _, printErr := fmt.Fprint(os.Stdout, string(out));
				printErr != nil {
				panic(printErr)
			}
		}
	}
}