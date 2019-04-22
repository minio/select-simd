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
	"io/ioutil"
	"testing"
)

func testCount(t *testing.T, f func(chunks []*bytes.Buffer, position uint64, condition string) (result int64)) {

	filename := "/Users/frankw/golang/src/github.com/fwessels/s3selectperf/data/parking-citations.csv"
	corpus, _ := ioutil.ReadFile(filename)

	const chunkSize = 0x10000

	chunks, _ := AlignChunks(chunkSize, corpus)

	testCases := []struct {
		make   string
		expected int64
	}{
		{"LAMB", 108},
		{"FERR", 3897},
		{"BENT", 11355},
		{"PORS", 134472},
		{"AUDI", 511914},
		{"BMW", 1205658},
		{"CHEV", 1799586},
		{"FORD", 2308959},
	}

	for _, tc := range testCases {
		result := f(chunks, 8, tc.make)
		if result != tc.expected {
			t.Errorf("Make test: got: %v want: %v", result, tc.expected)
		}
	}

	testCasesColor := []struct {
		color   string
		expected int64
	}{
		{"RE", 86037},
		{"GR", 56448},
		{"YE", 107007},
		{"WH", 301302},
		{"BL", 2052747},
	}

	for _, tc := range testCasesColor {
		result := f(chunks, 10, tc.color)
		if result != tc.expected {
			t.Errorf("Color test: got: %v want: %v", result, tc.expected)
		}
	}
}

func TestCount(t *testing.T) { testCount(t, ProcessCount) }
func TestCountParallel(t *testing.T) { testCount(t, ProcessCountParallel) }
