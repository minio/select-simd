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
	"bytes"
	"math/bits"
	"sync"
)

// count set bits in bitmask slice
func count64(bitMask []uint64) (result uint64) {

	for i := 0; i < len(bitMask); i++ {
		result += uint64(bits.OnesCount64(bitMask[i]))
	}
	return
}

// Perform select * operation
func ProcessSelectStar(chunks []*bytes.Buffer, position uint64, condition string) [][]byte {

	rows := bytes.Count(chunks[0].Bytes(), []byte{0x0a}) * 2

	rowIndices := make([]uint32, rows)
	maskCommas := make([]uint64, ((len(chunks[0].Bytes()) + 63) >> 6))
	colIndices := make([]uint32, rows)
	equal := make([]uint64, (rows+63)>>6)

	result := make([][]byte, 0, len(chunks))

	for _, chunk := range chunks {
		r := ParseCsvAdjusted(chunk.Bytes(), rowIndices, maskCommas)
		ExtractIndexForColumn(maskCommas, rowIndices[:r], colIndices[:r], position)
		EvaluateCompareString(chunk.Bytes(), colIndices[:r], condition, equal[:(r+63)>>6])
		bts := ExtractCsv(equal[:(r+63)>>6], chunk.Bytes(), rowIndices[:r])
		if len(bts) > 0 {
			result = append(result, bts)
		}
	}
	return result
}

// Perform select count(*) operation
func ProcessCount(chunks []*bytes.Buffer, position uint64, condition string) (result int64) {

	rows := bytes.Count(chunks[0].Bytes(), []byte{0x0a}) * 2

	rowIndices := make([]uint32, rows)
	maskCommas := make([]uint64, ((len(chunks[0].Bytes()) + 63) >> 6))
	colIndices := make([]uint32, rows)
	equal := make([]uint64, (rows+63)>>6)

	for _, chunk := range chunks {
		r := ParseCsvAdjusted(chunk.Bytes(), rowIndices, maskCommas)
		ExtractIndexForColumn(maskCommas, rowIndices[:r], colIndices[:r], position)
		EvaluateCompareString(chunk.Bytes(), colIndices[:r], condition, equal[:(r+63)>>6])
		result += int64(count64(equal[:(r+63)>>6]))
	}
	return
}

// Perform select count(*) operation in parallel
func ProcessCountParallel(chunks []*bytes.Buffer, position uint64, condition string) (result int64) {

	var wg sync.WaitGroup
	chunkChan := make(chan chunkInput)
	resultChan := make(chan chunkOutput)

	rows := bytes.Count(chunks[0].Bytes(), []byte{0x0a}) * 2

	// Start one go routine per CPU
	for i := 0; i < 16; /**cpu*/ i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processWorker(position, condition, rows, chunkChan, resultChan)
		}()
	}

	// Push chunks onto input channel
	go func() {
		for _, buf := range chunks {

			chunkChan <- chunkInput{buf: buf}
		}

		// Close input channel
		close(chunkChan)
	}()

	// Wait for workers to complete
	go func() {
		wg.Wait()
		close(resultChan) // Close output channel
	}()

	for r := range resultChan {
		result += int64(r.result)
	}

	return
}

type chunkInput struct {
	buf *bytes.Buffer
}

type chunkOutput struct {
	result uint64
}

// Worker routine to do parallel select
func processWorker(position uint64, condition string, rows int, chunkChan <-chan chunkInput, resultChan chan<- chunkOutput) {

	var rowIndices *[]uint32
	var maskCommas *[]uint64
	var colIndices *[]uint32
	var equal *[]uint64

	for c := range chunkChan {

		if rowIndices == nil {
			// TODO: Solve in a nicer manner
			_rowIndices := make([]uint32, rows)
			rowIndices = &_rowIndices
			_maskCommas := make([]uint64, ((len(c.buf.Bytes()) + 63) >> 6))
			maskCommas = &_maskCommas
			_colIndices := make([]uint32, rows)
			colIndices = &_colIndices
			_equal := make([]uint64, (rows+63)>>6)
			equal = &_equal
		}

		r := ParseCsvAdjusted(c.buf.Bytes(), *rowIndices, *maskCommas)
		ExtractIndexForColumn(*maskCommas, (*rowIndices)[:r], (*colIndices)[:r], position)
		EvaluateCompareString(c.buf.Bytes(), (*colIndices)[:r], condition, (*equal)[:(r+63)>>6])
		result := count64((*equal)[:(r+63)>>6])

		resultChan <- chunkOutput{result: result}
	}
}
