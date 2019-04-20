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
)

func testParse(t *testing.T, filename string) ([]uint32, []uint64) {

	rows := 8799051
	raw, _ := ioutil.ReadFile(filename)
	rowIndices := make([]uint32, rows)

	maskCommas := make([]uint64, ((len(raw) + 63) >> 6))

	ParseCsv(raw, rowIndices, maskCommas)

	return rowIndices, maskCommas
}

func TestParse(t *testing.T) {
	testParse(t, "/Users/frankw/golang/src/github.com/fwessels/s3selectperf/data/parking-citations.csv")

	// validate
}

func testExtract(t *testing.T, filename string) {

	rowIndices, maskCommas := testParse(t, filename)

	indices := make([]uint32, len(rowIndices))

	ExtractIndexForColumn(maskCommas, rowIndices, indices, 8)
}

func TestExtract(t *testing.T) {
	testExtract(t, "/Users/frankw/golang/src/github.com/fwessels/s3selectperf/data/parking-citations.csv")
}
