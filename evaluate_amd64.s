//+build !noasm,!appengine

//
// Minio Cloud Storage, (C) 2019 Minio, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// func evaluateCompareString(input []byte, indices []uint32, compare []byte, out []uint64)
TEXT ·evaluateCompareString(SB), 7, $0
	MOVQ in+0(FP), SI           // SI: &in
	MOVQ indices+24(FP), DI     // indices
	MOVQ indices_len+32(FP), R9 // length
	MOVQ compare+48(FP), AX     // &compare
	MOVQ out+72(FP), BX         // &out

	MOVOU (AX), X0            // comparison string
	MOVQ  trueMask+96(FP), AX

loop64:
	XORQ R10, R10
	MOVQ $1, R11

loop:
	MOVL      (DI), DX       // load index
	ADDQ      $4, DI
	MOVOU     (SI)(DX*1), X1 // load unaliged
	VPCMPEQB  X0, X1, X2
	VPMOVMSKB X2, CX
	ANDQ      AX, CX         // clear any top bits
	CMPQ      CX, AX
	JNE       skip
	ORQ       R11, R10       // set result bit

skip:
	SUBQ $1, R9
	JZ   done

	SHLQ $1, R11
	CMPQ R11, $0   // processed 64-bits?
	JNE  loop
	MOVQ R10, (BX) // store result
	ADDQ $8, BX
	JMP  loop64

done:
	MOVQ R10, (BX)
	RET
