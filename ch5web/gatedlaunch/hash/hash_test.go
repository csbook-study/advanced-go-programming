package hash

import (
	"fmt"
	"testing"
)

func BenchmarkMD5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		md5Hash()
	}
}

func BenchmarkSHA1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sha1Hash()
	}
}

func BenchmarkMurmurHash32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		murmur32()
	}
}

func BenchmarkMurmurHash64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		murmur64()
	}
}

func TestMurmur64distribute(t *testing.T) {
	t.Run("TestMurmur64distribute", func(t *testing.T) {
		var bucketSize = 10
		var bucketMap = map[uint64]int{}
		for i := 15000000000; i < 15000000000+10000000; i++ {
			hashInt := murmur64Str(fmt.Sprint(i)) % uint64(bucketSize)
			bucketMap[hashInt]++
		}
		fmt.Println("bucketMap: ", bucketMap)
	})
}
