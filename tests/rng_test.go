package tests

import (
	"testing"

	"github.com/DroidZed/go_lance/internal/utils"
)

func BenchmarkRNGSpeed(b *testing.B) {

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utils.RNG(59)
	}
}

func BenchmarkLinearRandomNumberGenerator(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utils.LinearRandomGenerator(89651649874945, 173, 17, 97, 3)
	}
}

func BenchmarkCombinedRNG(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		utils.LinearRandomGenerator(utils.RNG(89651649874945*int64(i)), 173, 17, 97, 3)
	}
}
