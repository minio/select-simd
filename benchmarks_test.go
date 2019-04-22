/*
 * Minio Cloud Storage, (C) 2019 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package selectsimd

import (
	"io/ioutil"
	"testing"
	"bytes"
)

func benchmark(b *testing.B, filename string) {
	raw, _ := ioutil.ReadFile(filename)
	rows := bytes.Count(raw, []byte{0x0a})

	rowIndices := make([]uint32, rows)
	maskCommas := make([]uint64, ((len(raw) + 63) >> 6))
	colIndices := make([]uint32, rows)
	equal := make([]uint64, (rows+63)>>6)

	b.SetBytes(int64(len(raw)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseCsv(raw, rowIndices, maskCommas)
		ExtractIndexForColumn(maskCommas, rowIndices, colIndices, 8)
		EvaluateCompareString(raw, colIndices, "HOND", equal)
	}
}

func BenchmarkSelectCsv(b *testing.B) {
	benchmark(b, "/Users/frankw/golang/src/github.com/fwessels/s3selectperf/data/parking-citations1x.csv")
}

func benchmarkParallel(b *testing.B, chunkSize int, filename string, cpus int) {

	corpus, _ := ioutil.ReadFile(filename)
	chunks, _ := AlignChunks(chunkSize, corpus)

	b.SetBytes(int64(len(corpus)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processCountParallel(chunks, 8, "PORS", cpus)
	}
}

func BenchmarkParallel_2cpus_256KB(b *testing.B) {
	benchmarkParallel(b, 0x40000, "/Users/frankw/golang/src/github.com/fwessels/s3selectperf/data/parking-citations1x.csv", 2)
}
func BenchmarkParallel_3cpus_256KB(b *testing.B) {
	benchmarkParallel(b, 0x40000, "/Users/frankw/golang/src/github.com/fwessels/s3selectperf/data/parking-citations1x.csv", 3)
}
func BenchmarkParallel_4cpus_256KB(b *testing.B) {
	benchmarkParallel(b, 0x40000, "/Users/frankw/golang/src/github.com/fwessels/s3selectperf/data/parking-citations1x.csv", 4)
}
