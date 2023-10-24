package cpu

func createOpcodesTable() [256]opcode {
	var opcodes []opcode

	// Create all the opcodes in any order

	// NOP
	opcodes = append(opcodes, createNOP(0x00))

	// LD rr,nn
	opcodes = append(opcodes, createLD_rr_nn(0x01, BC))
	opcodes = append(opcodes, createLD_rr_nn(0x11, DE))
	opcodes = append(opcodes, createLD_rr_nn(0x21, HL))
	opcodes = append(opcodes, createLD_rr_nn(0x31, SP))

	// LD (rr),A
	opcodes = append(opcodes, createLD_abs_rr_r(0x02, BC, A))
	opcodes = append(opcodes, createLD_abs_rr_r(0x12, DE, A))

	// LD (HL+),A
	opcodes = append(opcodes, createLD_inc_HL_A(0x22))
	// LD (HL-),A
	opcodes = append(opcodes, createLD_dec_HL_A(0x32))
	// LD A,(HL+)
	opcodes = append(opcodes, createLD_A_HL_inc(0x2A))
	// LD A,(HL-)
	opcodes = append(opcodes, createLD_A_HL_dec(0x3A))

	// INC rr
	opcodes = append(opcodes, createINC_rr(0x03, BC))
	opcodes = append(opcodes, createINC_rr(0x13, DE))
	opcodes = append(opcodes, createINC_rr(0x23, HL))
	opcodes = append(opcodes, createINC_rr(0x33, SP))

	// INC (HL)
	opcodes = append(opcodes, createINC_abs_HL(0x34))

	// DEC rr
	opcodes = append(opcodes, createDEC_rr(0x0B, BC))
	opcodes = append(opcodes, createDEC_rr(0x1B, DE))
	opcodes = append(opcodes, createDEC_rr(0x2B, HL))
	opcodes = append(opcodes, createDEC_rr(0x3B, SP))

	// DEC (HL)
	opcodes = append(opcodes, createDEC_abs_HL(0x35))

	// LD (HL), n
	opcodes = append(opcodes, createLD_abs_HL_n(0x36))

	// LD (nn),SP
	opcodes = append(opcodes, createLD_abs_nn_SP(0x08))

	// RLCA
	opcodes = append(opcodes, createRLCA(0x07))

	// RLA
	opcodes = append(opcodes, createRLA(0x17))

	// JR e
	opcodes = append(opcodes, createJR_e(0x18))

	// RRCA
	opcodes = append(opcodes, createRRCA(0x0F))

	// RRA
	opcodes = append(opcodes, createRRA(0x1F))

	// DAA
	opcodes = append(opcodes, createDAA(0x27))

	// INC r
	opcodes = append(opcodes, createINC_r(0x04, B))
	opcodes = append(opcodes, createINC_r(0x14, D))
	opcodes = append(opcodes, createINC_r(0x24, H))
	opcodes = append(opcodes, createINC_r(0x0C, C))
	opcodes = append(opcodes, createINC_r(0x1C, E))
	opcodes = append(opcodes, createINC_r(0x2C, L))
	opcodes = append(opcodes, createINC_r(0x3C, A))

	// DEC r
	opcodes = append(opcodes, createDEC_r(0x05, B))
	opcodes = append(opcodes, createDEC_r(0x15, D))
	opcodes = append(opcodes, createDEC_r(0x25, H))
	opcodes = append(opcodes, createDEC_r(0x0D, C))
	opcodes = append(opcodes, createDEC_r(0x1D, E))
	opcodes = append(opcodes, createDEC_r(0x2D, L))
	opcodes = append(opcodes, createDEC_r(0x3D, A))

	// LD r,n
	opcodes = append(opcodes, createLD_r_n(0x06, B))
	opcodes = append(opcodes, createLD_r_n(0x16, D))
	opcodes = append(opcodes, createLD_r_n(0x26, H))
	opcodes = append(opcodes, createLD_r_n(0x0E, C))
	opcodes = append(opcodes, createLD_r_n(0x1E, E))
	opcodes = append(opcodes, createLD_r_n(0x2E, L))
	opcodes = append(opcodes, createLD_r_n(0x3E, A))

	// ADD HL,rr
	opcodes = append(opcodes, createADD_HL_rr(0x09, BC))
	opcodes = append(opcodes, createADD_HL_rr(0x19, DE))
	opcodes = append(opcodes, createADD_HL_rr(0x29, HL))
	opcodes = append(opcodes, createADD_HL_rr(0x39, SP))

	// LD A,(rr)
	opcodes = append(opcodes, createLD_A_abs_rr(0x0A, BC))
	opcodes = append(opcodes, createLD_A_abs_rr(0x1A, DE))

	// DEC rr
	opcodes = append(opcodes, createDEC_rr(0x0B, BC))
	opcodes = append(opcodes, createDEC_rr(0x1B, DE))
	opcodes = append(opcodes, createDEC_rr(0x2B, HL))
	opcodes = append(opcodes, createDEC_rr(0x3B, SP))

	// JR cc,e
	opcodes = append(opcodes, createJR_cc_e(0x20, ZFlag, false))
	opcodes = append(opcodes, createJR_cc_e(0x30, CFlag, false))
	opcodes = append(opcodes, createJR_cc_e(0x28, ZFlag, true))
	opcodes = append(opcodes, createJR_cc_e(0x38, CFlag, true))

	// LD r,r
	opcodes = append(opcodes, createLD_r_r(0x40, B, B))
	opcodes = append(opcodes, createLD_r_r(0x41, B, C))
	opcodes = append(opcodes, createLD_r_r(0x42, B, D))
	opcodes = append(opcodes, createLD_r_r(0x43, B, E))
	opcodes = append(opcodes, createLD_r_r(0x44, B, H))
	opcodes = append(opcodes, createLD_r_r(0x45, B, L))
	opcodes = append(opcodes, createLD_r_r(0x47, B, A))

	opcodes = append(opcodes, createLD_r_r(0x48, C, B))
	opcodes = append(opcodes, createLD_r_r(0x49, C, C))
	opcodes = append(opcodes, createLD_r_r(0x4A, C, D))
	opcodes = append(opcodes, createLD_r_r(0x4B, C, E))
	opcodes = append(opcodes, createLD_r_r(0x4C, C, H))
	opcodes = append(opcodes, createLD_r_r(0x4D, C, L))
	opcodes = append(opcodes, createLD_r_r(0x4F, C, A))

	opcodes = append(opcodes, createLD_r_r(0x50, D, B))
	opcodes = append(opcodes, createLD_r_r(0x51, D, C))
	opcodes = append(opcodes, createLD_r_r(0x52, D, D))
	opcodes = append(opcodes, createLD_r_r(0x53, D, E))
	opcodes = append(opcodes, createLD_r_r(0x54, D, H))
	opcodes = append(opcodes, createLD_r_r(0x55, D, L))
	opcodes = append(opcodes, createLD_r_r(0x57, D, A))

	opcodes = append(opcodes, createLD_r_r(0x58, E, B))
	opcodes = append(opcodes, createLD_r_r(0x59, E, C))
	opcodes = append(opcodes, createLD_r_r(0x5A, E, D))
	opcodes = append(opcodes, createLD_r_r(0x5B, E, E))
	opcodes = append(opcodes, createLD_r_r(0x5C, E, H))
	opcodes = append(opcodes, createLD_r_r(0x5D, E, L))
	opcodes = append(opcodes, createLD_r_r(0x5F, E, A))

	opcodes = append(opcodes, createLD_r_r(0x60, H, B))
	opcodes = append(opcodes, createLD_r_r(0x61, H, C))
	opcodes = append(opcodes, createLD_r_r(0x62, H, D))
	opcodes = append(opcodes, createLD_r_r(0x63, H, E))
	opcodes = append(opcodes, createLD_r_r(0x64, H, H))
	opcodes = append(opcodes, createLD_r_r(0x65, H, L))
	opcodes = append(opcodes, createLD_r_r(0x67, H, A))

	opcodes = append(opcodes, createLD_r_r(0x68, L, B))
	opcodes = append(opcodes, createLD_r_r(0x69, L, C))
	opcodes = append(opcodes, createLD_r_r(0x6A, L, D))
	opcodes = append(opcodes, createLD_r_r(0x6B, L, E))
	opcodes = append(opcodes, createLD_r_r(0x6C, L, H))
	opcodes = append(opcodes, createLD_r_r(0x6D, L, L))
	opcodes = append(opcodes, createLD_r_r(0x6F, L, A))

	opcodes = append(opcodes, createLD_r_r(0x78, A, B))
	opcodes = append(opcodes, createLD_r_r(0x79, A, C))
	opcodes = append(opcodes, createLD_r_r(0x7A, A, D))
	opcodes = append(opcodes, createLD_r_r(0x7B, A, E))
	opcodes = append(opcodes, createLD_r_r(0x7C, A, H))
	opcodes = append(opcodes, createLD_r_r(0x7D, A, L))
	opcodes = append(opcodes, createLD_r_r(0x7F, A, A))

	// LD r,(HL)
	opcodes = append(opcodes, createLD_r_abs_HL(0x46, B))
	opcodes = append(opcodes, createLD_r_abs_HL(0x4E, C))
	opcodes = append(opcodes, createLD_r_abs_HL(0x56, D))
	opcodes = append(opcodes, createLD_r_abs_HL(0x5E, E))
	opcodes = append(opcodes, createLD_r_abs_HL(0x66, H))
	opcodes = append(opcodes, createLD_r_abs_HL(0x6E, L))
	opcodes = append(opcodes, createLD_r_abs_HL(0x7E, A))

	// LD (HL),r
	opcodes = append(opcodes, createLD_abs_rr_r(0x70, HL, B))
	opcodes = append(opcodes, createLD_abs_rr_r(0x71, HL, C))
	opcodes = append(opcodes, createLD_abs_rr_r(0x72, HL, D))
	opcodes = append(opcodes, createLD_abs_rr_r(0x73, HL, E))
	opcodes = append(opcodes, createLD_abs_rr_r(0x74, HL, H))
	opcodes = append(opcodes, createLD_abs_rr_r(0x75, HL, L))
	opcodes = append(opcodes, createLD_abs_rr_r(0x77, HL, A))

	// ADD r
	opcodes = append(opcodes, createADD_r(0x80, B))
	opcodes = append(opcodes, createADD_r(0x81, C))
	opcodes = append(opcodes, createADD_r(0x82, D))
	opcodes = append(opcodes, createADD_r(0x83, E))
	opcodes = append(opcodes, createADD_r(0x84, H))
	opcodes = append(opcodes, createADD_r(0x85, L))
	opcodes = append(opcodes, createADD_abs_HL(0x86))
	opcodes = append(opcodes, createADD_r(0x87, A))

	// ADC r
	opcodes = append(opcodes, createADC_r(0x88, B))
	opcodes = append(opcodes, createADC_r(0x89, C))
	opcodes = append(opcodes, createADC_r(0x8A, D))
	opcodes = append(opcodes, createADC_r(0x8B, E))
	opcodes = append(opcodes, createADC_r(0x8C, H))
	opcodes = append(opcodes, createADC_r(0x8D, L))
	opcodes = append(opcodes, createADC_abs_HL(0x8E))
	opcodes = append(opcodes, createADC_r(0x8F, A))

	// SUB r
	opcodes = append(opcodes, createSUB_r(0x90, B))
	opcodes = append(opcodes, createSUB_r(0x91, C))
	opcodes = append(opcodes, createSUB_r(0x92, D))
	opcodes = append(opcodes, createSUB_r(0x93, E))
	opcodes = append(opcodes, createSUB_r(0x94, H))
	opcodes = append(opcodes, createSUB_r(0x95, L))
	opcodes = append(opcodes, createSUB_abs_HL(0x96))
	opcodes = append(opcodes, createSUB_r(0x97, A))

	// SBC r
	opcodes = append(opcodes, createSBC_r(0x98, B))
	opcodes = append(opcodes, createSBC_r(0x99, C))
	opcodes = append(opcodes, createSBC_r(0x9A, D))
	opcodes = append(opcodes, createSBC_r(0x9B, E))
	opcodes = append(opcodes, createSBC_r(0x9C, H))
	opcodes = append(opcodes, createSBC_r(0x9D, L))
	opcodes = append(opcodes, createSBC_abs_HL(0x9E))
	opcodes = append(opcodes, createSBC_r(0x9F, A))

	// AND r
	opcodes = append(opcodes, createAND_r(0xA0, B))
	opcodes = append(opcodes, createAND_r(0xA1, C))
	opcodes = append(opcodes, createAND_r(0xA2, D))
	opcodes = append(opcodes, createAND_r(0xA3, E))
	opcodes = append(opcodes, createAND_r(0xA4, H))
	opcodes = append(opcodes, createAND_r(0xA5, L))
	opcodes = append(opcodes, createAND_abs_HL(0xA6))
	opcodes = append(opcodes, createAND_r(0xA7, A))

	// XOR r
	opcodes = append(opcodes, createXOR_r(0xA8, B))
	opcodes = append(opcodes, createXOR_r(0xA9, C))
	opcodes = append(opcodes, createXOR_r(0xAA, D))
	opcodes = append(opcodes, createXOR_r(0xAB, E))
	opcodes = append(opcodes, createXOR_r(0xAC, H))
	opcodes = append(opcodes, createXOR_r(0xAD, L))
	opcodes = append(opcodes, createXOR_abs_HL(0xAE))
	opcodes = append(opcodes, createXOR_r(0xAF, A))

	// OR r
	opcodes = append(opcodes, createOR_r(0xB0, B))
	opcodes = append(opcodes, createOR_r(0xB1, C))
	opcodes = append(opcodes, createOR_r(0xB2, D))
	opcodes = append(opcodes, createOR_r(0xB3, E))
	opcodes = append(opcodes, createOR_r(0xB4, H))
	opcodes = append(opcodes, createOR_r(0xB5, L))
	opcodes = append(opcodes, createOR_abs_HL(0xB6))
	opcodes = append(opcodes, createOR_r(0xB7, A))

	// CP r
	opcodes = append(opcodes, createCP_r(0xB8, B))
	opcodes = append(opcodes, createCP_r(0xB9, C))
	opcodes = append(opcodes, createCP_r(0xBA, D))
	opcodes = append(opcodes, createCP_r(0xBB, E))
	opcodes = append(opcodes, createCP_r(0xBC, H))
	opcodes = append(opcodes, createCP_r(0xBD, L))
	opcodes = append(opcodes, createCP_abs_HL(0xBE))
	opcodes = append(opcodes, createCP_r(0xBF, A))

	// RET cc
	opcodes = append(opcodes, createRET_cc(0xC0, ZFlag, false))
	opcodes = append(opcodes, createRET_cc(0xD0, CFlag, false))
	opcodes = append(opcodes, createRET_cc(0xC8, ZFlag, true))
	opcodes = append(opcodes, createRET_cc(0xD8, CFlag, true))

	// POP rr
	opcodes = append(opcodes, createPOP_rr(0xC1, BC))
	opcodes = append(opcodes, createPOP_rr(0xD1, DE))
	opcodes = append(opcodes, createPOP_rr(0xE1, HL))
	opcodes = append(opcodes, createPOP_rr(0xF1, AF))

	// JP cc,nn
	opcodes = append(opcodes, createJP_cc_nn(0xC2, ZFlag, false))
	opcodes = append(opcodes, createJP_cc_nn(0xD2, CFlag, false))
	opcodes = append(opcodes, createJP_cc_nn(0xCA, ZFlag, true))
	opcodes = append(opcodes, createJP_cc_nn(0xDA, CFlag, true))

	// JP nn
	opcodes = append(opcodes, createJP_nn(0xC3))

	// CALL cc,nn
	opcodes = append(opcodes, createCALL_cc_nn(0xC4, ZFlag, false))
	opcodes = append(opcodes, createCALL_cc_nn(0xD4, CFlag, false))
	opcodes = append(opcodes, createCALL_cc_nn(0xCC, ZFlag, true))
	opcodes = append(opcodes, createCALL_cc_nn(0xDC, CFlag, true))

	// PUSH rr
	opcodes = append(opcodes, createPUSH_rr(0xC5, BC))
	opcodes = append(opcodes, createPUSH_rr(0xD5, DE))
	opcodes = append(opcodes, createPUSH_rr(0xE5, HL))
	opcodes = append(opcodes, createPUSH_rr(0xF5, AF))

	// ADD n
	opcodes = append(opcodes, createADD_n(0xC6))
	// SUB n
	opcodes = append(opcodes, createSUB_n(0xD6))
	// AND n
	opcodes = append(opcodes, createAND_n(0xE6))
	// OR n
	opcodes = append(opcodes, createOR_n(0xF6))

	// RST n
	opcodes = append(opcodes, createRST_n(0xC7, 0x00))
	opcodes = append(opcodes, createRST_n(0xD7, 0x10))
	opcodes = append(opcodes, createRST_n(0xE7, 0x20))
	opcodes = append(opcodes, createRST_n(0xF7, 0x30))
	opcodes = append(opcodes, createRST_n(0xCF, 0x08))
	opcodes = append(opcodes, createRST_n(0xDF, 0x18))
	opcodes = append(opcodes, createRST_n(0xEF, 0x28))
	opcodes = append(opcodes, createRST_n(0xFF, 0x38))

	// RET
	opcodes = append(opcodes, createRET(0xC9))

	// RETI
	opcodes = append(opcodes, createRETI(0xD9))

	// CALL nn
	opcodes = append(opcodes, createCALL_nn(0xCD))

	// ADC n
	opcodes = append(opcodes, createADC_n(0xCE))
	// SBC n
	opcodes = append(opcodes, createSBC_n(0xDE))
	// XOR n
	opcodes = append(opcodes, createXOR_n(0xEE))
	// CP n
	opcodes = append(opcodes, createCP_n(0xFE))

	// LDH (n),A
	opcodes = append(opcodes, createLDH_imed_n_A(0xE0))
	// LDH A,(n)
	opcodes = append(opcodes, createLDH_A_imed_n(0xF0))

	// DI
	opcodes = append(opcodes, createDI(0xF3))

	// LDH (C),A
	opcodes = append(opcodes, createLDH_imed_C_A(0xE2))
	// LDH A,(C)
	opcodes = append(opcodes, createLDH_A_imed_C(0xF2))

	// EI
	opcodes = append(opcodes, createEI(0xFB))

	// JP HL
	opcodes = append(opcodes, createJP_HL(0xE9))

	// CCF, SCF, CPL
	opcodes = append(opcodes, createCCF(0x3F))
	opcodes = append(opcodes, createSCF(0x37))
	opcodes = append(opcodes, createCPL(0x2F))

	// LD (nn),A & LD A,(nn)
	opcodes = append(opcodes, createLD_nn_A(0xEA))
	opcodes = append(opcodes, createLD_A_nn(0xFA))

	// LD HL,SP+e
	opcodes = append(opcodes, createLD_HL_SP_plus_e(0xF8))

	// Populate the opcodes table according to the opcode value set on each opcode
	var table [256]opcode
	for _, x := range opcodes {
		table[x.opcode()] = x
	}
	return table
}

func createCBOpcodesTable() [256]opcode {
	var opcodes []opcode

	// RL r
	opcodes = append(opcodes, createCB_RL_r(0x10, B))
	opcodes = append(opcodes, createCB_RL_r(0x11, C))
	opcodes = append(opcodes, createCB_RL_r(0x12, D))
	opcodes = append(opcodes, createCB_RL_r(0x13, E))
	opcodes = append(opcodes, createCB_RL_r(0x14, H))
	opcodes = append(opcodes, createCB_RL_r(0x15, L))
	opcodes = append(opcodes, createCB_RL_r(0x17, A))

	// RR
	opcodes = append(opcodes, createCB_RR_r(0x18, B))
	opcodes = append(opcodes, createCB_RR_r(0x19, C))
	opcodes = append(opcodes, createCB_RR_r(0x1A, D))
	opcodes = append(opcodes, createCB_RR_r(0x1B, E))
	opcodes = append(opcodes, createCB_RR_r(0x1C, H))
	opcodes = append(opcodes, createCB_RR_r(0x1D, L))
	opcodes = append(opcodes, createCB_RR_r(0x1F, A))

	// SLA r
	opcodes = append(opcodes, createCB_SLA_r(0x20, B))
	opcodes = append(opcodes, createCB_SLA_r(0x21, C))
	opcodes = append(opcodes, createCB_SLA_r(0x22, D))
	opcodes = append(opcodes, createCB_SLA_r(0x23, E))
	opcodes = append(opcodes, createCB_SLA_r(0x24, H))
	opcodes = append(opcodes, createCB_SLA_r(0x25, L))
	opcodes = append(opcodes, createCB_SLA_r(0x27, A))

	// SWAP
	opcodes = append(opcodes, createCB_SWAP_r(0x30, B))
	opcodes = append(opcodes, createCB_SWAP_r(0x31, C))
	opcodes = append(opcodes, createCB_SWAP_r(0x32, D))
	opcodes = append(opcodes, createCB_SWAP_r(0x33, E))
	opcodes = append(opcodes, createCB_SWAP_r(0x34, H))
	opcodes = append(opcodes, createCB_SWAP_r(0x35, L))
	opcodes = append(opcodes, createCB_SWAP_r(0x37, A))

	// SRL
	opcodes = append(opcodes, createCB_SRL_r(0x38, B))
	opcodes = append(opcodes, createCB_SRL_r(0x39, C))
	opcodes = append(opcodes, createCB_SRL_r(0x3A, D))
	opcodes = append(opcodes, createCB_SRL_r(0x3B, E))
	opcodes = append(opcodes, createCB_SRL_r(0x3C, H))
	opcodes = append(opcodes, createCB_SRL_r(0x3D, L))
	opcodes = append(opcodes, createCB_SRL_r(0x3F, A))

	// BIT b,r
	opcodes = append(opcodes, createCB_BIT_b_r(0x40, 0, B))
	opcodes = append(opcodes, createCB_BIT_b_r(0x41, 0, C))
	opcodes = append(opcodes, createCB_BIT_b_r(0x42, 0, D))
	opcodes = append(opcodes, createCB_BIT_b_r(0x43, 0, E))
	opcodes = append(opcodes, createCB_BIT_b_r(0x44, 0, H))
	opcodes = append(opcodes, createCB_BIT_b_r(0x45, 0, L))
	opcodes = append(opcodes, createCB_BIT_b_HL(0x46, 0))
	opcodes = append(opcodes, createCB_BIT_b_r(0x47, 0, A))

	opcodes = append(opcodes, createCB_BIT_b_r(0x48, 1, B))
	opcodes = append(opcodes, createCB_BIT_b_r(0x49, 1, C))
	opcodes = append(opcodes, createCB_BIT_b_r(0x4A, 1, D))
	opcodes = append(opcodes, createCB_BIT_b_r(0x4B, 1, E))
	opcodes = append(opcodes, createCB_BIT_b_r(0x4C, 1, H))
	opcodes = append(opcodes, createCB_BIT_b_r(0x4D, 1, L))
	opcodes = append(opcodes, createCB_BIT_b_HL(0x4E, 1))
	opcodes = append(opcodes, createCB_BIT_b_r(0x4F, 1, A))

	opcodes = append(opcodes, createCB_BIT_b_r(0x50, 2, B))
	opcodes = append(opcodes, createCB_BIT_b_r(0x51, 2, C))
	opcodes = append(opcodes, createCB_BIT_b_r(0x52, 2, D))
	opcodes = append(opcodes, createCB_BIT_b_r(0x53, 2, E))
	opcodes = append(opcodes, createCB_BIT_b_r(0x54, 2, H))
	opcodes = append(opcodes, createCB_BIT_b_r(0x55, 2, L))
	opcodes = append(opcodes, createCB_BIT_b_HL(0x56, 2))
	opcodes = append(opcodes, createCB_BIT_b_r(0x57, 2, A))

	opcodes = append(opcodes, createCB_BIT_b_r(0x58, 3, B))
	opcodes = append(opcodes, createCB_BIT_b_r(0x59, 3, C))
	opcodes = append(opcodes, createCB_BIT_b_r(0x5A, 3, D))
	opcodes = append(opcodes, createCB_BIT_b_r(0x5B, 3, E))
	opcodes = append(opcodes, createCB_BIT_b_r(0x5C, 3, H))
	opcodes = append(opcodes, createCB_BIT_b_r(0x5D, 3, L))
	opcodes = append(opcodes, createCB_BIT_b_HL(0x5E, 3))
	opcodes = append(opcodes, createCB_BIT_b_r(0x5F, 3, A))

	opcodes = append(opcodes, createCB_BIT_b_r(0x60, 4, B))
	opcodes = append(opcodes, createCB_BIT_b_r(0x61, 4, C))
	opcodes = append(opcodes, createCB_BIT_b_r(0x62, 4, D))
	opcodes = append(opcodes, createCB_BIT_b_r(0x63, 4, E))
	opcodes = append(opcodes, createCB_BIT_b_r(0x64, 4, H))
	opcodes = append(opcodes, createCB_BIT_b_r(0x65, 4, L))
	opcodes = append(opcodes, createCB_BIT_b_HL(0x66, 4))
	opcodes = append(opcodes, createCB_BIT_b_r(0x67, 4, A))

	opcodes = append(opcodes, createCB_BIT_b_r(0x68, 5, B))
	opcodes = append(opcodes, createCB_BIT_b_r(0x69, 5, C))
	opcodes = append(opcodes, createCB_BIT_b_r(0x6A, 5, D))
	opcodes = append(opcodes, createCB_BIT_b_r(0x6B, 5, E))
	opcodes = append(opcodes, createCB_BIT_b_r(0x6C, 5, H))
	opcodes = append(opcodes, createCB_BIT_b_r(0x6D, 5, L))
	opcodes = append(opcodes, createCB_BIT_b_HL(0x6E, 5))
	opcodes = append(opcodes, createCB_BIT_b_r(0x6F, 5, A))

	opcodes = append(opcodes, createCB_BIT_b_r(0x70, 6, B))
	opcodes = append(opcodes, createCB_BIT_b_r(0x71, 6, C))
	opcodes = append(opcodes, createCB_BIT_b_r(0x72, 6, D))
	opcodes = append(opcodes, createCB_BIT_b_r(0x73, 6, E))
	opcodes = append(opcodes, createCB_BIT_b_r(0x74, 6, H))
	opcodes = append(opcodes, createCB_BIT_b_r(0x75, 6, L))
	opcodes = append(opcodes, createCB_BIT_b_HL(0x76, 6))
	opcodes = append(opcodes, createCB_BIT_b_r(0x77, 6, A))

	opcodes = append(opcodes, createCB_BIT_b_r(0x78, 7, B))
	opcodes = append(opcodes, createCB_BIT_b_r(0x79, 7, C))
	opcodes = append(opcodes, createCB_BIT_b_r(0x7A, 7, D))
	opcodes = append(opcodes, createCB_BIT_b_r(0x7B, 7, E))
	opcodes = append(opcodes, createCB_BIT_b_r(0x7C, 7, H))
	opcodes = append(opcodes, createCB_BIT_b_r(0x7D, 7, L))
	opcodes = append(opcodes, createCB_BIT_b_HL(0x7E, 7))
	opcodes = append(opcodes, createCB_BIT_b_r(0x7F, 7, A))

	// RES b,r
	opcodes = append(opcodes, createCB_RES_b_r(0x80, 0, B))
	opcodes = append(opcodes, createCB_RES_b_r(0x81, 0, C))
	opcodes = append(opcodes, createCB_RES_b_r(0x82, 0, D))
	opcodes = append(opcodes, createCB_RES_b_r(0x83, 0, E))
	opcodes = append(opcodes, createCB_RES_b_r(0x84, 0, H))
	opcodes = append(opcodes, createCB_RES_b_r(0x85, 0, L))
	opcodes = append(opcodes, createCB_RES_b_HL(0x86, 0))
	opcodes = append(opcodes, createCB_RES_b_r(0x87, 0, A))

	opcodes = append(opcodes, createCB_RES_b_r(0x88, 1, B))
	opcodes = append(opcodes, createCB_RES_b_r(0x89, 1, C))
	opcodes = append(opcodes, createCB_RES_b_r(0x8A, 1, D))
	opcodes = append(opcodes, createCB_RES_b_r(0x8B, 1, E))
	opcodes = append(opcodes, createCB_RES_b_r(0x8C, 1, H))
	opcodes = append(opcodes, createCB_RES_b_r(0x8D, 1, L))
	opcodes = append(opcodes, createCB_RES_b_HL(0x8E, 1))
	opcodes = append(opcodes, createCB_RES_b_r(0x8F, 1, A))

	opcodes = append(opcodes, createCB_RES_b_r(0x90, 2, B))
	opcodes = append(opcodes, createCB_RES_b_r(0x91, 2, C))
	opcodes = append(opcodes, createCB_RES_b_r(0x92, 2, D))
	opcodes = append(opcodes, createCB_RES_b_r(0x93, 2, E))
	opcodes = append(opcodes, createCB_RES_b_r(0x94, 2, H))
	opcodes = append(opcodes, createCB_RES_b_r(0x95, 2, L))
	opcodes = append(opcodes, createCB_RES_b_HL(0x96, 2))
	opcodes = append(opcodes, createCB_RES_b_r(0x97, 2, A))

	opcodes = append(opcodes, createCB_RES_b_r(0x98, 3, B))
	opcodes = append(opcodes, createCB_RES_b_r(0x99, 3, C))
	opcodes = append(opcodes, createCB_RES_b_r(0x9A, 3, D))
	opcodes = append(opcodes, createCB_RES_b_r(0x9B, 3, E))
	opcodes = append(opcodes, createCB_RES_b_r(0x9C, 3, H))
	opcodes = append(opcodes, createCB_RES_b_r(0x9D, 3, L))
	opcodes = append(opcodes, createCB_RES_b_HL(0x9E, 3))
	opcodes = append(opcodes, createCB_RES_b_r(0x9F, 3, A))

	opcodes = append(opcodes, createCB_RES_b_r(0xA0, 4, B))
	opcodes = append(opcodes, createCB_RES_b_r(0xA1, 4, C))
	opcodes = append(opcodes, createCB_RES_b_r(0xA2, 4, D))
	opcodes = append(opcodes, createCB_RES_b_r(0xA3, 4, E))
	opcodes = append(opcodes, createCB_RES_b_r(0xA4, 4, H))
	opcodes = append(opcodes, createCB_RES_b_r(0xA5, 4, L))
	opcodes = append(opcodes, createCB_RES_b_HL(0xA6, 4))
	opcodes = append(opcodes, createCB_RES_b_r(0xA7, 4, A))

	opcodes = append(opcodes, createCB_RES_b_r(0xA8, 5, B))
	opcodes = append(opcodes, createCB_RES_b_r(0xA9, 5, C))
	opcodes = append(opcodes, createCB_RES_b_r(0xAA, 5, D))
	opcodes = append(opcodes, createCB_RES_b_r(0xAB, 5, E))
	opcodes = append(opcodes, createCB_RES_b_r(0xAC, 5, H))
	opcodes = append(opcodes, createCB_RES_b_r(0xAD, 5, L))
	opcodes = append(opcodes, createCB_RES_b_HL(0xAE, 5))
	opcodes = append(opcodes, createCB_RES_b_r(0xAF, 5, A))

	opcodes = append(opcodes, createCB_RES_b_r(0xB0, 6, B))
	opcodes = append(opcodes, createCB_RES_b_r(0xB1, 6, C))
	opcodes = append(opcodes, createCB_RES_b_r(0xB2, 6, D))
	opcodes = append(opcodes, createCB_RES_b_r(0xB3, 6, E))
	opcodes = append(opcodes, createCB_RES_b_r(0xB4, 6, H))
	opcodes = append(opcodes, createCB_RES_b_r(0xB5, 6, L))
	opcodes = append(opcodes, createCB_RES_b_HL(0xB6, 6))
	opcodes = append(opcodes, createCB_RES_b_r(0xB7, 6, A))

	opcodes = append(opcodes, createCB_RES_b_r(0xB8, 7, B))
	opcodes = append(opcodes, createCB_RES_b_r(0xB9, 7, C))
	opcodes = append(opcodes, createCB_RES_b_r(0xBA, 7, D))
	opcodes = append(opcodes, createCB_RES_b_r(0xBB, 7, E))
	opcodes = append(opcodes, createCB_RES_b_r(0xBC, 7, H))
	opcodes = append(opcodes, createCB_RES_b_r(0xBD, 7, L))
	opcodes = append(opcodes, createCB_RES_b_HL(0xBE, 7))
	opcodes = append(opcodes, createCB_RES_b_r(0xBF, 7, A))

	// Populate the opcodes table according to the opcode value set on each opcode
	var table [256]opcode
	for _, x := range opcodes {
		table[x.opcode()] = x
	}
	return table
}
