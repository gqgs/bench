package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	defaultIterations = 5
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		println("usage: ./bench [cmd] [args]")
		return
	}
	iterations := defaultIterations
	if envIterations, err := strconv.Atoi(os.Getenv("BENCH_ITERATIONS")); err == nil {
		iterations = envIterations
	}

	diff := make([]time.Duration, 0, iterations)
	var cmdArgs []string
	if len(args) > 1 {
		cmdArgs = args[1:]
	}
	for iteration := iterations; iteration > 0; iteration-- {
		cmd := exec.Command(args[0], cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		start := time.Now()
		if err := cmd.Run(); err != nil {
			return
		}
		delta := time.Now().Sub(start)
		diff = append(diff, delta)
	}

	var sum time.Duration
	for i := range diff {
		sum += diff[i]
	}
	mean := sum / time.Duration(iterations)
	var variance time.Duration
	for i := range diff {
		variance += (diff[i] - mean) * (diff[i] - mean)
	}
	variance /= time.Duration(iterations)
	stdDeviation := time.Duration(math.Sqrt(float64(variance)))

	fmt.Printf("\nbenchmark stats:\n iterations: %d\n mean: %v\n std deviation: %v\n", iterations, mean, stdDeviation)
}
