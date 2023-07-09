package main

import (
	"testing"
)

func BenchmarkParseCron(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := ParseCron("* * * * *")
		_ = err
	}
}
