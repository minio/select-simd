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
	"encoding/hex"
	"log"
	"strconv"
	"strings"
)

// Parse the output of hex.Dump([]byte) back into a byte slice
func dump2hex(data string) []byte {
	addr, addrFrom := uint64(0), ^uint64(0)
	blob := make([]byte, 0)
	lines := strings.Split(data, "\n")
	for _, l := range lines {
		sections := strings.Split(l, "  ")
		if len(sections) < 1 || len(sections[0]) == 0 {
			continue
		} else if sections[0] == "*" {
			addrFrom = addr
			continue
		}

		if a, err := strconv.ParseUint("0x"+sections[0], 0, 64); err != nil {
			log.Fatal(err)
		} else {
			addr = a
			if addrFrom != ^uint64(0) {
				for a := addrFrom + 16; a < addr; a += 16 {
					blob = append(blob, blob[len(blob)-16:len(blob)]...)
				}
				addrFrom = ^uint64(0)
			}
		}

		if len(sections) < 2 {
			continue
		}
		for s := 1; s <= 2; s++ {
			parts := strings.Split(sections[s], " ")
			decoded, err := hex.DecodeString(strings.Join(parts, ""))
			if err != nil {
				log.Fatal(err)
			}
			blob = append(blob, decoded...)
		}
	}
	return blob
}
