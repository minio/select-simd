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
)

func AlignChunks(chunkSize int, corpus []byte) (chunks []*bytes.Buffer, nextChunk *bytes.Buffer) {

	chunk := &bytes.Buffer{}
	chunk.Grow(chunkSize)
	for start := 0; start < len(corpus); {

		end := start + chunkSize - chunk.Len()
		if end <= len(corpus) {
			chunk.Write(corpus[start:end]) // write full chunk
		} else {
			chunk.Write(corpus[start:])                             // write remaining part of corpus
			chunk.Write(bytes.Repeat([]byte{0x0}, end-len(corpus))) // fill rest with zeros
		}
		start = end

		// compute trailingBytes for this chunk
		trailingBytes := chunkSize - (bytes.LastIndex(chunk.Bytes(), []byte{0x0a}) + 1)

		nextChunk = &bytes.Buffer{}
		nextChunk.Grow(chunkSize)

		if trailingBytes > 0 {
			nextChunk.Write(chunk.Bytes()[chunkSize-trailingBytes:])

			// zero out trailing bytes (not strictly necessary)
			copy(chunk.Bytes()[chunkSize-trailingBytes:], bytes.Repeat([]byte{0x0}, trailingBytes))
		}

		chunks = append(chunks, chunk)
		chunk = nextChunk
	}

	return
}
