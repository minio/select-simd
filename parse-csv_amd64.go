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

//go:noescape
func scanForDelimiterAndSeparator(raw []byte, indices []uint32, mask []uint64, delimiter uint64, separator uint64) (rows uint64)

//go:noescape
func seekIndexForPosition(mask []uint64, rowIndices []uint32, indices []uint32, position uint64)
