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

func ParseCsv(raw []byte, rowIndices []uint32, maskSeparator []uint64) (rows uint64) {
	return scanForDelimiterAndSeparator(raw, rowIndices, maskSeparator, 0x0a, 0x2c)
}

func ParseCsvAdjusted(raw []byte, rowIndices []uint32, maskSeparator []uint64) (rows uint64) {
	rows = scanForDelimiterAndSeparator(raw, rowIndices, maskSeparator, 0x0a, 0x2c)
	if rows > 0 {
		// TODO: Decide on inclusiveness of last delimiter or not
		rows -= 1
	}
	return rows
}

func ExtractIndexForColumn(maskSeparator []uint64, rowIndices []uint32, indices []uint32, position uint64) {

	if position > 0 {
		seekIndexForPosition(maskSeparator, rowIndices, indices, position)
	} else {
		index := uint32(0)
		for i, c := range rowIndices {
			indices[i] = index
			index = c + 1
		}
	}
}

func ExtractCsv(bitMask []uint64, raw []byte, rowIndices []uint32) (result []byte) {

	rows := count64(bitMask)

	avgPerRow := uint64(len(raw) / len(rowIndices))
	buf := make([]byte, 0, rows*avgPerRow*5/4)
	if bitMask[0]&1 != 0 {
		buf = append(buf, raw[0:rowIndices[0]-1]...)
		buf = append(buf, 0x0)
	}

	for b := uint64(1); b < 64; b++ {
		if bitMask[0]&(1<<b) != 0 {
			buf = append(buf, raw[rowIndices[b-1]+1:rowIndices[b]-1]...)
			buf = append(buf, 0x0)
		}
	}

	for i := 1; i < len(bitMask)-1; i++ {
		for b := uint64(0); b < 64; b++ {
			if bitMask[i]&(1<<b) != 0 {
				index := uint64(i)*64 + b
				buf = append(buf, raw[rowIndices[index-1]+1:rowIndices[index]-1]...)
				buf = append(buf, 0x0)
			}
		}
	}

	lastMask := bitMask[len(bitMask)-1]
	for b := uint64(0); b < 64; b++ {
		if lastMask&(1<<b) != 0 {
			index := uint64(len(bitMask)-1)*64 + b
			//if rowIndices[index-1]+1 >= rowIndices[index]-1 {
			//	// TODO: Check why this boundary condition may occur
			//	fmt.Println("*** rowIndices[index-1]+1 smaller than rowIndices[index]-1")
			//	break
			//}

			buf = append(buf, raw[rowIndices[index-1]+1:rowIndices[index]-1]...)
			buf = append(buf, 0x0)
		}
	}

	return buf
}
