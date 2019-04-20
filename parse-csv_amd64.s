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

// func scanForDelimiterAndSeparator(raw []byte, indices []uint32, mask []uint64, delimiter uint64, separator uint64) (rows uint64)
TEXT ·scanForDelimiterAndSeparator(SB), 7, $0
	MOVQ   raw+0(FP), SI      // SI: &raw
 	MOVQ   raw_len+8(FP), R9  // R9: len(raw)
	MOVQ   indices+24(FP), DI // DI: &indices
	MOVQ   mask+48(FP), DX    // DX: &mask

	// TODO: Load indices_len and make sure we do not write beyond
	// TODO: Load mask_len and make sure we do not write beyond

	SHRQ   $6, R9             // len(in) / 64
	CMPQ   R9, $0
	JEQ    done

	MOVQ   delimiter+72(FP), AX  // get newline
    MOVQ   AX, X0
  	VPBROADCASTB X0, Y0

	MOVQ   separator+80(FP), AX  // get separator
    MOVQ   AX, X6
  	VPBROADCASTB X6, Y6

    XORQ   BX, BX

loop:
    // First scan for separator
	VPCMPEQB   0x00(SI)(BX*1), Y6,  Y7
	VPCMPEQB   0x20(SI)(BX*1), Y6,  Y8
    VPMOVMSKB  Y7,  AX
    VPMOVMSKB  Y8,  CX
    SHLQ       $32, CX
	ORQ        CX,  AX
    MOVQ       AX, (DX)
	ADDQ       $8, DX         // mask += 8

    // Scan for delimiter
	VPCMPEQB   0x00(SI)(BX*1), Y0,  Y1
	VPCMPEQB   0x20(SI)(BX*1), Y0,  Y2
    VPMOVMSKB  Y1,  AX
    VPMOVMSKB  Y2,  CX
    SHLQ       $32, CX
	ORQ        CX,  AX
    JZ         skipCtz

loopCtz:
    TZCNTQ AX, R10
    ADDQ   $4, DI
    ADDQ   BX, R10
    BLSRQ  AX, AX
    MOVL   R10, -4(DI)
    JNZ    loopCtz

skipCtz:
 	ADDQ   $64, BX
	SUBQ   $1, R9
	JNZ    loop

done:
	MOVQ   indices+24(FP), SI // reload indices pointer
    SUBQ   SI, DI
    ADDQ   $4, DI             // make final pointer inclusive
    SHRQ   $2, DI
    MOVQ   DI, rows+88(FP)    // store result
    VZEROUPPER
    RET

// func seekIndexForPosition(mask []uint64, rowIndices []uint32, indices []uint32, position uint64)
TEXT ·seekIndexForPosition(SB), 7, $0
	MOVQ   mask+0(FP), SI             // SI: &mask
	MOVQ   mask_len+8(FP), R9         // R9: len(mask)
	MOVQ   rowIndices+24(FP), R12     // R12: &rowIndices
	MOVQ   rowIndices_len+32(FP), R13 // R13: len(rowIndices)
	MOVQ   indices+48(FP), R10        // R10: &indices

    XORQ   CX, CX
    JMP    skipFirst

rowLoop:
    MOVL   (R12), CX
    ADDQ   $1, CX      // add one extra character to skip the newline
    ADDQ   $4, R12

skipFirst:
	MOVQ   position+72(FP), R11      // R11: position

    // Initialize offset into mask
    MOVQ   CX, DX
    SHRQ   $6, DX
    SHLQ   $3, DX

    // Setup mask to clear bits from previous row
    MOVQ   $1, AX
    ANDQ   $0x3f, CX
    SHLQ   CX, AX
    DECQ   AX
    NOTQ   AX

    MOVQ   (SI)(DX*1), BX  // load initial mask for this row
    ANDQ   AX, BX          // clear bits from previous row

    XORQ   AX, AX
    JMP    loopCtz

doneCtz:
	ADDQ   $8, DX
    MOVQ   (SI)(DX*1), BX // load next mask to examine

loopCtz:
    TZCNTQ BX, AX       // count number of trailing zero bits
    CMPQ   AX, $64
    JEQ    doneCtz
    BLSRQ  BX, BX       // clear trailing bit
    SUBQ   $1, R11      // decrement position we are seeking
    JNZ    loopCtz

done:
    SHLQ   $3, DX
    ADDQ   DX, AX
    ADDQ   $1, AX       // skip beyond last separator
    MOVL   AX, (R10)
    ADDQ   $4, R10

    SUBQ   $1, R13
    JNZ    rowLoop
    RET
