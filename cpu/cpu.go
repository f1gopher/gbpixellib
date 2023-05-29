package cpu

import (
	"fmt"
	"go-boy/log"
	"go-boy/memory"
	"os"
)

type register int

const (
	AF register = iota
	BC
	DE
	HL
	SP
	PC
	A
	F
	B
	C
	D
	E
	H
	L
)

func (r register) String() string {
	return [...]string{"AF", "BC", "DE", "HL", "SP", "PC", "A", "F", "B", "C", "D", "E", "H", "L"}[r]
}

var opcodeCycles = [...]int{
	1, // NOP",
	3, // LB BC,nn",
	1, // LD (BC),A",
	1, // INC BC",
	1, // INC B",
	1, // DEC B",
	2, // LD B,n",
	1, // RLCA",
	5, // LD (nn),SP",
	1, // ADD HL,BC",
	2, // LD A,(BC)",
	1, // DEC BC",
	1, // INC C",
	1, // DEC C",
	2, // LD C,n",
	1, // RRCA",
	1, // STOP",
	3, // LD DE,nn",
	1, // LD (DE),A",
	1, // INC DE",
	1, // INC D",
	1, // DEC D",
	2, // LD D,n",
	1, // RLA",
	3, // JR e",
	1, // ADD HL,DE",
	2, // LD A,(DE)",
	1, // DEC DE",
	1, // INC E",
	1, // DEC E",
	2, // LD E,n",
	1, // RRA",
	1, // JR NZ,e",
	1, // LD HL,nn",
	1, // LD (HL+),A",
	1, // INC HL",
	1, // INC H",
	1, // DEC H",
	1, // LD H,n",
	1, // DAA",
	1, // JR Z, e",
	1, // ADD HL,HL",
	1, // LD A,(HL+)",
	1, // DEC HL",
	1, // INC L",
	1, // DEC L",
	1, // LD L,n",
	1, // CPL",
	1, // JR NC,e",
	1, // LD SP,nn",
	1, // LD (HL-),A",
	1, // INC SP",
	1, // INC (HL)",
	1, // DEC (HL)",
	3, // LD (HL),n",
	1, // SCF",
	1, // JR C,e",
	1, // ADD HL,SP",
	1, // LD A,(HL-)",
	1, // DEC SP",
	1, // INC A",
	1, // DEC A",
	1, // LD A,n",
	1, // CCF",
	1, // LD B,B",
	1, // LD B,C",
	1, // LD B,D",
	1, // LD B,E",
	1, // LD B,H",
	1, // LD B,L",
	1, // LD B,(HL)",
	1, // LD B,A",
	1, // LD C,B",
	1, // LD C,C",
	1, // LD C,D",
	1, // LD C,E",
	1, // LD C,H",
	1, // LD C,L",
	1, // LD C,(HL)",
	1, // LD C,A",
	1, // LD D,B",
	1, // LD D,C",
	1, // LD D,D",
	1, // LD D,E",
	1, // LD D,H",
	1, // LD D,L",
	1, // LD D,(HL)",
	1, // LD D,A",
	1, // LD E,B",
	1, // LD E,C",
	1, // LD E,D",
	1, // LD E,E",
	1, // LD E,H",
	1, // LD E,L",
	1, // LD E,(HL)",
	1, // LD E,A",
	1, // LD H,B",
	1, // LD H,C",
	1, // LD H,D",
	1, // LD H,E",
	1, // LD H,H",
	1, // LD H,L",
	1, // LD H,(HL)",
	1, // LD H,A",
	1, // LD L,B",
	1, // LD L,C",
	1, // LD L,D",
	1, // LD L,E",
	1, // LD L,H",
	1, // LD L,L",
	1, // LD L,(HL)",
	1, // LD L,A",
	1, // LD (HL),B",
	1, // LD (HL),C",
	1, // LD (HL),D",
	1, // LD (HL),E",
	1, // LD (HL),H",
	1, // LD (HL),L",
	1, // HALT",
	1, // LD (HL),A",
	1, // LD A,B",
	1, // LD A,C",
	1, // LD A,D",
	1, // LD A,E",
	1, // LD A,H",
	1, // LD A,L",
	1, // LD A,(HL)",
	1, // LD A,A",
	1, // ADD B",
	1, // ADD C",
	1, // ADD D",
	1, // ADD E",
	1, // ADD H",
	1, // ADD L",
	1, // ADD (HL)",
	1, // ADD A",
	1, // ADC B",
	1, // ADC C",
	1, // ADC D",
	1, // ADC E",
	1, // ADC H",
	1, // ADC L",
	1, // ADC (HL)",
	1, // ADC A",
	1, // SUB B",
	1, // SUB C",
	1, // SUB D",
	1, // SUB E",
	1, // SUB H",
	1, // SUB L",
	1, // SUB (HL)",
	1, // SUB A",
	1, // SBC B",
	1, // SBC C",
	1, // SBC D",
	1, // SBC E",
	1, // SBC H",
	1, // SBC L",
	1, // SBC (HL)",
	1, // SBC A",
	1, // AND B",
	1, // AND C",
	1, // AND D",
	1, // AND E",
	1, // AND H",
	1, // AND L",
	1, // AND (HL)",
	1, // AND A",
	1, // XOR B",
	1, // XOR C",
	1, // XOR D",
	1, // XOR E",
	1, // XOR H",
	1, // XOR L",
	1, // XOR (HL)",
	1, // XOR A",
	1, // OR B",
	1, // OR C",
	1, // OR D",
	1, // OR E",
	1, // OR H",
	1, // OR L",
	1, // OR (HL)",
	1, // OR A",
	1, // CP B",
	1, // CP C",
	1, // CP D",
	1, // CP E",
	1, // CP H",
	1, // CP L",
	1, // CP (HL)",
	1, // CP A",
	1, // RET NZ",
	1, // POP BC",
	1, // JP NZ,nn",
	1, // JP nn",
	1, // CALL NZ,nn",
	1, // PUSH BC",
	1, // ADD n",
	1, // RST 0x00",
	1, // RET Z",
	1, // RET",
	1, // JP Z,nn",
	1, // CB",
	1, // CALL Z,nn",
	1, // CALL nn",
	1, // ADC n",
	1, // RST 0x08",
	1, // RET NC",
	1, // POP DE",
	1, // JP NC,nn",
	1, // <NULL>",
	1, // CALL NC,nn",
	1, // PUSH DE",
	1, // SUB n",
	1, // RST 0x10",
	1, // RET C",
	1, // RETI",
	1, // JP C,nn",
	1, // <NULL>",
	1, // CALL C,nn",
	1, // <NULL>",
	1, // SBC n",
	1, // RST 0x18",
	1, // LDH (n),A",
	1, // POP HL",
	1, // LDH (C),A",
	1, // <NULL>",
	1, // <NULL>",
	1, // PUSH HL",
	1, // AND n",
	1, // RST 0x20",
	1, // ADD SP,e",
	1, // JP HL",
	1, // LD (nn),A",
	1, // <NULL>",
	1, // <NULL>",
	1, // <NULL>",
	1, // XOR n",
	1, // RST 0x28",
	1, // LDH A,(n)",
	1, // POP AF",
	1, // LDH A,(C)",
	1, // DI",
	1, // <NULL>",
	1, // PUSH AF",
	1, // OR n",
	1, // RST 0x30",
	1, // LD HL,SP+e",
	1, // LD SP,HL",
	1, // LD A,(nn)",
	1, // EI",
	1, // <NULL>",
	1, // <NULL>",
	1, // CP n",
	1, // RST 0x38",
}

var opcodeNames = [...]string{
	"NOP",
	"LB BC,nn",
	"LD (BC),A",
	"INC BC",
	"INC B",
	"DEC B",
	"LD B,n",
	"RLCA",
	"LD (nn),SP",
	"ADD HL,BC",
	"LD A,(BC)",
	"DEC BC",
	"INC C",
	"DEC C",
	"LD C,n",
	"RRCA",
	"STOP",
	"LD DE,nn",
	"LD (DE),A",
	"INC DE",
	"INC D",
	"DEC D",
	"LD D,n",
	"RLA",
	"JR e",
	"ADD HL,DE",
	"LD A,(DE)",
	"DEC DE",
	"INC E",
	"DEC E",
	"LD E,n",
	"RRA",
	"JR NZ,e",
	"LD HL,nn",
	"LD (HL+),A",
	"INC HL",
	"INC H",
	"DEC H",
	"LD H,n",
	"DAA",
	"JR Z, e",
	"ADD HL,HL",
	"LD A,(HL+)",
	"DEC HL",
	"INC L",
	"DEC L",
	"LD L,n",
	"CPL",
	"JR NC,e",
	"LD SP,nn",
	"LD (HL-),A",
	"INC SP",
	"INC (HL)",
	"DEC (HL)",
	"LD (HL),n",
	"SCF",
	"JR C,e",
	"ADD HL,SP",
	"LD A,(HL-)",
	"DEC SP",
	"INC A",
	"DEC A",
	"LD A,n",
	"CCF",
	"LD B,B",
	"LD B,C",
	"LD B,D",
	"LD B,E",
	"LD B,H",
	"LD B,L",
	"LD B,(HL)",
	"LD B,A",
	"LD C,B",
	"LD C,C",
	"LD C,D",
	"LD C,E",
	"LD C,H",
	"LD C,L",
	"LD C,(HL)",
	"LD C,A",
	"LD D,B",
	"LD D,C",
	"LD D,D",
	"LD D,E",
	"LD D,H",
	"LD D,L",
	"LD D,(HL)",
	"LD D,A",
	"LD E,B",
	"LD E,C",
	"LD E,D",
	"LD E,E",
	"LD E,H",
	"LD E,L",
	"LD E,(HL)",
	"LD E,A",
	"LD H,B",
	"LD H,C",
	"LD H,D",
	"LD H,E",
	"LD H,H",
	"LD H,L",
	"LD H,(HL)",
	"LD H,A",
	"LD L,B",
	"LD L,C",
	"LD L,D",
	"LD L,E",
	"LD L,H",
	"LD L,L",
	"LD L,(HL)",
	"LD L,A",
	"LD (HL),B",
	"LD (HL),C",
	"LD (HL),D",
	"LD (HL),E",
	"LD (HL),H",
	"LD (HL),L",
	"HALT",
	"LD (HL),A",
	"LD A,B",
	"LD A,C",
	"LD A,D",
	"LD A,E",
	"LD A,H",
	"LD A,L",
	"LD A,(HL)",
	"LD A,A",
	"ADD B",
	"ADD C",
	"ADD D",
	"ADD E",
	"ADD H",
	"ADD L",
	"ADD (HL)",
	"ADD A",
	"ADC B",
	"ADC C",
	"ADC D",
	"ADC E",
	"ADC H",
	"ADC L",
	"ADC (HL)",
	"ADC A",
	"SUB B",
	"SUB C",
	"SUB D",
	"SUB E",
	"SUB H",
	"SUB L",
	"SUB (HL)",
	"SUB A",
	"SBC B",
	"SBC C",
	"SBC D",
	"SBC E",
	"SBC H",
	"SBC L",
	"SBC (HL)",
	"SBC A",
	"AND B",
	"AND C",
	"AND D",
	"AND E",
	"AND H",
	"AND L",
	"AND (HL)",
	"AND A",
	"XOR B",
	"XOR C",
	"XOR D",
	"XOR E",
	"XOR H",
	"XOR L",
	"XOR (HL)",
	"XOR A",
	"OR B",
	"OR C",
	"OR D",
	"OR E",
	"OR H",
	"OR L",
	"OR (HL)",
	"OR A",
	"CP B",
	"CP C",
	"CP D",
	"CP E",
	"CP H",
	"CP L",
	"CP (HL)",
	"CP A",
	"RET NZ",
	"POP BC",
	"JP NZ,nn",
	"JP nn",
	"CALL NZ,nn",
	"PUSH BC",
	"ADD n",
	"RST 0x00",
	"RET Z",
	"RET",
	"JP Z,nn",
	"CB",
	"CALL Z,nn",
	"CALL nn",
	"ADC n",
	"RST 0x08",
	"RET NC",
	"POP DE",
	"JP NC,nn",
	"<NULL>",
	"CALL NC,nn",
	"PUSH DE",
	"SUB n",
	"RST 0x10",
	"RET C",
	"RETI",
	"JP C,nn",
	"<NULL>",
	"CALL C,nn",
	"<NULL>",
	"SBC n",
	"RST 0x18",
	"LDH (n),A",
	"POP HL",
	"LDH (C),A",
	"<NULL>",
	"<NULL>",
	"PUSH HL",
	"AND n",
	"RST 0x20",
	"ADD SP,e",
	"JP HL",
	"LD (nn),A",
	"<NULL>",
	"<NULL>",
	"<NULL>",
	"XOR n",
	"RST 0x28",
	"LDH A,(n)",
	"POP AF",
	"LDH A,(C)",
	"DI",
	"<NULL>",
	"PUSH AF",
	"OR n",
	"RST 0x30",
	"LD HL,SP+e",
	"LD SP,HL",
	"LD A,(nn)",
	"EI",
	"<NULL>",
	"<NULL>",
	"CP n",
	"RST 0x38",
}

var cbNames = [...]string{
	"RLC B",
	"RLC C",
	"RLC D",
	"RLC E",
	"RLC H",
	"RLC L",
	"RLC (HL)",
	"RLC A",
	"RRC B",
	"RRC C",
	"RRC D",
	"RRC E",
	"RRC H",
	"RRC L",
	"RRC (HL)",
	"RRC A",
	"RL B",
	"RL C",
	"RL D",
	"RL E",
	"RL H",
	"RL L",
	"RL (HL)",
	"RL A",
	"RR B",
	"RR C",
	"RR D",
	"RR E",
	"RR H",
	"RR L",
	"RR (HL)",
	"RR A",
	"SLA B",
	"SLA C",
	"SLA D",
	"SLA E",
	"SLA H",
	"SLA L",
	"SLA (HL)",
	"SLA A",
	"SRA B",
	"SRA C",
	"SRA D",
	"SRA E",
	"SRA H",
	"SRA L",
	"SRA (HL)",
	"SRA A",
	"SWAP B",
	"SWAP C",
	"SWAP D",
	"SWAP E",
	"SWAP H",
	"SWAP L",
	"SWAP (HL)",
	"SWAP A",
	"SRL B",
	"SRL C",
	"SRL D",
	"SRL E",
	"SRL H",
	"SRL L",
	"SRL (HL)",
	"SRL A",
	"BIT 0, B",
	"BIT 0,C",
	"BIT 0,D",
	"BIT 0,E",
	"BIT 0,H",
	"BIT 0,L",
	"BIT 0,(HL)",
	"BIT 0,A",
	"BIT 1,B",
	"BIT 1,C",
	"BIT 1,D",
	"BIT 1,E",
	"BIT 1,H",
	"BIT 1,L",
	"BIT 1,(HL)",
	"BIT 1,A",
	"BIT 2,B",
	"BIT 2,C",
	"BIT 2,D",
	"BIT 2,E",
	"BIT 2,H",
	"BIT 2,L",
	"BIT 2,(HL)",
	"BIT 2,A",
	"BIT 3,B",
	"BIT 3,C",
	"BIT 3,D",
	"BIT 3,E",
	"BIT 3,H",
	"BIT 3,L",
	"BIT 3,(HL)",
	"BIT 3,A",
	"BIT 4,B",
	"BIT 4,C",
	"BIT 4,D",
	"BIT 4,E",
	"BIT 4,H",
	"BIT 4,L",
	"BIT 4,(HL)",
	"BIT 4,A",
	"BIT 5,B",
	"BIT 5,C",
	"BIT 5,D",
	"BIT 5,E",
	"BIT 5,H",
	"BIT 5,L",
	"BIT 5,(HL)",
	"BIT 5,A",
	"BIT 6,B",
	"BIT 6,C",
	"BIT 6,D",
	"BIT 6,E",
	"BIT 6,H",
	"BIT 6,L",
	"BIT 6,(HL)",
	"BIT 6,A",
	"BIT 7,B",
	"BIT 7,C",
	"BIT 7,D",
	"BIT 7,E",
	"BIT 7,H",
	"BIT 7,L",
	"BIT 7,(HL)",
	"BIT 7,A",
	"RES 0,B",
	"RES 0,C",
	"RES 0,D",
	"RES 0,E",
	"RES 0,H",
	"RES 0,L",
	"RES 0,(HL)",
	"RES 0,A",
	"RES 1,B",
	"RES 1,C",
	"RES 1,D",
	"RES 1,E",
	"RES 1,H",
	"RES 1,L",
	"RES 1,(HL)",
	"RES 1,A",
	"RES 2,B",
	"RES 2,C",
	"RES 2,D",
	"RES 2,E",
	"RES 2,H",
	"RES 2,L",
	"RES 2,(HL)",
	"RES 3,B",
	"RES 3,C",
	"RES 3,D",
	"RES 3,E",
	"RES 3,H",
	"RES 3,L",
	"RES 3,(HL)",
	"RES 3,A",
	"RES 4,B",
	"RES 4,C",
	"RES 4,D",
	"RES 4,E",
	"RES 4,H",
	"RES 4,L",
	"RES 4,(HL)",
	"RES 4,A",
	"RES 5,B",
	"RES 5,C",
	"RES 5,D",
	"RES 5,E",
	"RES 5,H",
	"RES 5,L",
	"RES 5,(HL)",
	"RES 5,A",
	"RES 6,B",
	"RES 6,C",
	"RES 6,D",
	"RES 6,E",
	"RES 6,H",
	"RES 6,L",
	"RES 6,HL",
	"RES 6,A",
	"RES 7,B",
	"RES 7,C",
	"RES 7,D",
	"RES 7,E",
	"RES 7,H",
	"RES 7,L",
	"RES 7,(HL)",
	"RES 7,A",
	"SET 0,C",
	"SET 0,D",
	"SET 0,E",
	"SET 0,H",
	"SET 0,L",
	"SET 0,(HL)",
	"SET 0,A",
	"SET 1,B",
	"SET 1,C",
	"SET 1,D",
	"SET 1,E",
	"SET 1,H",
	"SET 1,L",
	"SET 1,(HL)",
	"SET 1,A",
	"SET 2,B",
	"SET 2,C",
	"SET 2,D",
	"SET 2,E",
	"SET 2,H",
	"SET 2,L",
	"SET 2,(HL)",
	"SET 3,B",
	"SET 3,C",
	"SET 3,D",
	"SET 3,E",
	"SET 3,H",
	"SET 3,L",
	"SET 3,(HL)",
	"SET 3,A",
	"SET 4,B",
	"SET 4,C",
	"SET 4,D",
	"SET 4,E",
	"SET 4,H",
	"SET 4,L",
	"SET 4,(HL)",
	"SET 4,A",
	"SET 5,B",
	"SET 5,C",
	"SET 5,D",
	"SET 5,E",
	"SET 5,H",
	"SET 5,L",
	"SET 5,(HL)",
	"SET 5,A",
	"SET 6,B",
	"SET 6,C",
	"SET 6,D",
	"SET 6,E",
	"SET 6,H",
	"SET 6,L",
	"SET 6,HL",
	"SET 6,A",
	"SET 7,B",
	"SET 7,C",
	"SET 7,D",
	"SET 7,E",
	"SET 7,H",
	"SET 7,L",
	"SET 7,(HL)",
	"SET 7,A",
}

type CPU struct {
	memory *memory.Memory
	regAF  uint16
	regBC  uint16
	regDE  uint16
	regHL  uint16
	regSP  uint16
	regPC  uint16

	opcodes     map[byte]func()
	opcodesUsed [255]bool

	//debugLogFile *os.File

	log *log.Log

	count int
}

func CreateCPU(log *log.Log, memory *memory.Memory) *CPU {
	chip := CPU{
		log:    log,
		memory: memory,
		count:  0,
	}

	// TODO - opcodes range from 0-256 so we can use an array instead
	chip.opcodes = map[byte]func(){
		0x40: chip.op_LD_B_B,
		0x41: chip.op_LD_B_C,
		0x42: chip.op_LD_B_D,
		0x43: chip.op_LD_B_E,
		0x44: chip.op_LD_B_H,
		0x45: chip.op_LD_B_L,
		0x47: chip.op_LD_B_A,
		0x48: chip.op_LD_C_B,
		0x49: chip.op_LD_C_C,
		0x4A: chip.op_LD_C_D,
		0x4B: chip.op_LD_C_E,
		0x4C: chip.op_LD_C_H,
		0x4D: chip.op_LD_C_L,
		0x4F: chip.op_LD_C_A,
		0x50: chip.op_LD_D_B,
		0x51: chip.op_LD_D_C,
		0x52: chip.op_LD_D_D,
		0x53: chip.op_LD_D_E,
		0x54: chip.op_LD_D_H,
		0x55: chip.op_LD_D_L,
		0x57: chip.op_LD_D_A,
		0x58: chip.op_LD_E_B,
		0x59: chip.op_LD_E_C,
		0x5A: chip.op_LD_E_D,
		0x5B: chip.op_LD_E_E,
		0x5C: chip.op_LD_E_H,
		0x5D: chip.op_LD_E_L,
		0x5F: chip.op_LD_E_A,
		0x60: chip.op_LD_H_B,
		0x61: chip.op_LD_H_C,
		0x62: chip.op_LD_H_D,
		0x63: chip.op_LD_H_E,
		0x64: chip.op_LD_H_H,
		0x65: chip.op_LD_H_L,
		0x67: chip.op_LD_H_A,
		0x68: chip.op_LD_L_B,
		0x69: chip.op_LD_L_C,
		0x6A: chip.op_LD_L_D,
		0x6B: chip.op_LD_L_E,
		0x6C: chip.op_LD_L_H,
		0x6D: chip.op_LD_L_L,
		0x6F: chip.op_LD_L_A,
		0x78: chip.op_LD_A_B,
		0x79: chip.op_LD_A_C,
		0x7A: chip.op_LD_A_D,
		0x7B: chip.op_LD_A_E,
		0x7C: chip.op_LD_A_H,
		0x7D: chip.op_LD_A_L,
		0x7F: chip.op_LD_A_A,
		0x00: chip.op_NOP,
		0xC3: chip.op_JP,
		0x01: chip.op_LD_BC_nn,
		0x11: chip.op_LD_DE_nn,
		0x21: chip.op_LD_HL_nn,
		0x31: chip.op_LD_SP_nn,
		0x2A: chip.op_LD_A_HL_plus,
		0x12: chip.op_LD_DE_A,
		0x0C: chip.op_INC_C,
		0x1C: chip.op_INC_E,
		0x2C: chip.op_INC_L,
		0x3C: chip.op_INC_A,
		0x04: chip.op_INC_B,
		0x14: chip.op_INC_D,
		0x24: chip.op_INC_H,
		0x0D: chip.op_DEC_C,
		0x1D: chip.op_DEC_E,
		0x2D: chip.op_DEC_L,
		0x3D: chip.op_DEC_A,
		0x05: chip.op_DEC_B,
		0x15: chip.op_DEC_D,
		0x25: chip.op_DEC_H,
		0x03: chip.op_INC_BC,
		0x13: chip.op_INC_DE,
		0x23: chip.op_INC_HL,
		0x33: chip.op_INC_HL,
		0x0B: chip.op_DEC_BC,
		0x1B: chip.op_DEC_DE,
		0x2B: chip.op_DEC_HL,
		0x3B: chip.op_DEC_SP,
		0x20: chip.op_JR_NZ_e,
		0x18: chip.op_JR_e,
		0x30: chip.op_JR_NC_e,
		0x28: chip.op_JR_Z_e,
		0x38: chip.op_JR_C_e,
		0xF3: chip.op_DI,
		0xEA: chip.op_LD_nn_A,
		0x06: chip.op_LD_B_n,
		0x16: chip.op_LD_D_n,
		0x26: chip.op_LD_H_n,
		0x0E: chip.op_LD_C_n,
		0x1E: chip.op_LD_E_n,
		0x2E: chip.op_LD_L_n,
		0x3E: chip.op_LD_A_n,
		0xE0: chip.op_LDH_n_A,
		0xCD: chip.op_CALL_nn,
		0xC9: chip.op_RET,
		0xC5: chip.op_PUSH_BC,
		0xD5: chip.op_PUSH_DE,
		0xE5: chip.op_PUSH_HL,
		0xF5: chip.op_PUSH_AF,
		0xC1: chip.op_POP_BC,
		0xD1: chip.op_POP_DE,
		0xE1: chip.op_POP_HL,
		0xF1: chip.op_POP_AF,
		0xB0: chip.op_OR_B,
		0xB1: chip.op_OR_C,
		0xB2: chip.op_OR_D,
		0xB3: chip.op_OR_E,
		0xB4: chip.op_OR_H,
		0xB5: chip.op_OR_L,
		0xB7: chip.op_OR_A,
		0xA0: chip.op_AND_B,
		0xA1: chip.op_AND_C,
		0xA2: chip.op_AND_D,
		0xA3: chip.op_AND_E,
		0xA4: chip.op_AND_H,
		0xA5: chip.op_AND_L,
		0xA7: chip.op_AND_A,
		0xA8: chip.op_XOR_B,
		0xA9: chip.op_XOR_C,
		0xAA: chip.op_XOR_D,
		0xAB: chip.op_XOR_E,
		0xAC: chip.op_XOR_H,
		0xAD: chip.op_XOR_L,
		0xAF: chip.op_XOR_A,
		0xF0: chip.op_LDH_A_n,
		0xFE: chip.op_CP_n,
		0xFA: chip.op_LD_A_nn,
		0xE6: chip.op_AND_n,
		0xC4: chip.op_CALL_NZ_nn,
		0xCC: chip.op_CALL_Z_nn,
		0xD4: chip.op_CALL_NC_nn,
		0xDC: chip.op_CALL_C_nn,
		0x70: chip.op_LD_HL_B,
		0x71: chip.op_LD_HL_C,
		0x72: chip.op_LD_HL_D,
		0x73: chip.op_LD_HL_E,
		0x74: chip.op_LD_HL_H,
		0x75: chip.op_LD_HL_L,
		0x77: chip.op_LD_HL_A,
		0x1A: chip.op_LD_A_DE,
		0x22: chip.op_LD_HL_plus_A,
		0x32: chip.op_LD_HL_sub_A,
		0xC6: chip.op_ADD_n,
		0xD6: chip.op_SUB_n,
		0x46: chip.op_LD_B_HL,
		0x56: chip.op_LD_D_HL,
		0x66: chip.op_LD_H_HL,
		0x4E: chip.op_LD_C_HL,
		0x5E: chip.op_LD_E_HL,
		0x6E: chip.op_LD_L_HL,
		0x7E: chip.op_LD_A_HL,
		//0xAE: chip.op_XOR_HL,
		0xCB: chip.op_CB_op,
		0xE2: chip.op_LDH_C_A,
		0x17: chip.op_RL_A,
		0x90: chip.op_SUB_B,
		0x91: chip.op_SUB_C,
		0x92: chip.op_SUB_D,
		0x93: chip.op_SUB_E,
		0x94: chip.op_SUB_H,
		0x95: chip.op_SUB_L,
		0x96: chip.op_SUB_HL,
		0x97: chip.op_SUB_A,
		0xB8: chip.op_CP_B,
		0xB9: chip.op_CP_C,
		0xBA: chip.op_CP_D,
		0xBB: chip.op_CP_E,
		0xBC: chip.op_CP_H,
		0xBD: chip.op_CP_L,
		0xBE: chip.op_CP_HL,
		0xBF: chip.op_CP_A,
		0x80: chip.op_ADD_B,
		0x81: chip.op_ADD_C,
		0x82: chip.op_ADD_D,
		0x83: chip.op_ADD_E,
		0x84: chip.op_ADD_H,
		0x85: chip.op_ADD_L,
		0x86: chip.op_ADD_HL,
		0x87: chip.op_ADD_A,
		0x36: chip.op_LD_HL_n,
	}

	return &chip
}

func (c *CPU) InitForTestROM() {
	c.setRegByte(A, 0x01)
	c.setRegByte(F, 0xB0)
	c.setRegByte(B, 0x00)
	c.setRegByte(C, 0x13)
	c.setRegByte(D, 0x00)
	c.setRegByte(E, 0xD8)
	c.setRegByte(H, 0x01)
	c.setRegByte(L, 0x4D)
	c.setRegShort(SP, 0xFFFE)
	c.setRegShort(PC, 0x0100)
}

func (c *CPU) Init() {
	c.setRegShort(PC, 0x0000) //0x100)
	c.setRegShort(AF, 0x01B0)
	c.setRegShort(BC, 0x0013)
	c.setRegShort(DE, 0x00D8)
	c.setRegShort(HL, 0x014D)
	c.setRegShort(SP, 0xFFFE)
	c.memory.WriteByte(0xFF05, 0x00)
	c.memory.WriteByte(0xFF06, 0x00)
	c.memory.WriteByte(0xFF07, 0x00)
	c.memory.WriteByte(0xFF10, 0x80)
	c.memory.WriteByte(0xFF11, 0xBF)
	c.memory.WriteByte(0xFF12, 0xF3)
	c.memory.WriteByte(0xFF14, 0xBF)
	c.memory.WriteByte(0xFF16, 0x3F)
	c.memory.WriteByte(0xFF17, 0x00)
	c.memory.WriteByte(0xFF19, 0xBF)
	c.memory.WriteByte(0xFF1A, 0x7F)
	c.memory.WriteByte(0xFF1B, 0xFF)
	c.memory.WriteByte(0xFF1C, 0x9F)
	c.memory.WriteByte(0xFF1E, 0xBF)
	c.memory.WriteByte(0xFF20, 0xFF)
	c.memory.WriteByte(0xFF21, 0x00)
	c.memory.WriteByte(0xFF22, 0x00)
	c.memory.WriteByte(0xFF23, 0xBF)
	c.memory.WriteByte(0xFF24, 0x77)
	c.memory.WriteByte(0xFF25, 0xF3)
	c.memory.WriteByte(0xFF26, 0xF1)
	c.memory.WriteByte(0xFF40, 0x91)
	c.memory.WriteByte(0xFF42, 0x00)
	c.memory.WriteByte(0xFF43, 0x00)
	c.memory.WriteByte(0xFF45, 0x00)
	c.memory.WriteByte(0xFF47, 0xFC)
	c.memory.WriteByte(0xFF48, 0xFF)
	c.memory.WriteByte(0xFF49, 0xFF)
	c.memory.WriteByte(0xFF4A, 0x00)
	c.memory.WriteByte(0xFF4B, 0x00)
	c.memory.WriteByte(0xFFFF, 0x00)
}

//func (c *CPU) SetDebugLog(file string) error {
//	var err error
//	c.debugLogFile, err = os.Create(file)
//
//	if err != nil {
//		return errors.Join(errors.New("Unable to set debug log"), err)
//	}
//
//	return nil
//}

func (c *CPU) PushAndReplacePC(newPC uint16) {

	currentPC := c.getRegShort(PC)

	c.regSP -= 2
	c.memory.WriteShort(c.regSP, currentPC)

	c.setRegShort(PC, newPC)
}

func (c *CPU) Tick() int {
	c.debugLog()

	opcode := c.memory.ReadByte(c.regPC)

	//fmt.Printf("%d - Op: 0x%02X ", c.count, opcode)

	c.opcodesUsed[opcode] = true

	executor := c.getOpcode(opcode)

	// TODO - is this the best place to increment the program counter?
	c.regPC++
	executor()

	c.count++

	// TODO - return correct cycles count for opcode
	return opcodeCycles[opcode]
}

func (c *CPU) getOpcode(opcode byte) func() {

	executor, exists := c.opcodes[opcode]
	if !exists {
		panic(fmt.Sprintf("Unsupported opcode: 0x%02X %s", opcode, opcodeNames[opcode]))
	}

	name := opcodeNames[opcode]
	if opcode == 0xCB {
		name = name + " - " + cbNames[0]
	}

	return executor
}

func (c *CPU) getRegShort(reg register) uint16 {
	switch reg {
	case AF:
		return c.regAF
	case BC:
		return c.regBC
	case DE:
		return c.regDE
	case HL:
		return c.regHL
	case SP:
		return c.regSP
	case PC:
		return c.regPC
	default:
		panic(fmt.Sprintf("Invalid register for get short : %s", reg.String()))
	}
}

func (c *CPU) getRegByte(reg register) byte {
	switch reg {
	case A:
		return getHighByte(c.regAF)
	case F:
		return getLowByte(c.regAF)
	case B:
		return getHighByte(c.regBC)
	case C:
		return getLowByte(c.regBC)
	case D:
		return getHighByte(c.regDE)
	case E:
		return getLowByte(c.regDE)
	case H:
		return getHighByte(c.regHL)
	case L:
		return getLowByte(c.regHL)
	default:
		panic(fmt.Sprintf("Invalid register for byte: %s", reg.String()))
	}
}

func (c *CPU) getFlagZ() bool { return c.getRegBit(F, 7) }
func (c *CPU) getFlagN() bool { return c.getRegBit(F, 6) }
func (c *CPU) getFlagH() bool { return c.getRegBit(F, 5) }
func (c *CPU) getFlagC() bool { return c.getRegBit(F, 4) }

func (c *CPU) getRegBit(reg register, bit int) bool {
	return memory.GetBit(c.getRegByte(reg), bit)
}

func (c *CPU) setRegShort(reg register, value uint16) {
	switch reg {
	case AF:
		c.regAF = value
	case BC:
		c.regBC = value
	case DE:
		c.regDE = value
	case HL:
		c.regHL = value
	case SP:
		c.regSP = value
	case PC:
		c.regPC = value
	default:
		panic(fmt.Sprintf("Invalid register for set short: %s", reg.String()))
	}
}

func (c *CPU) setRegByte(reg register, value byte) {
	switch reg {
	case A:
		setHighByte(&c.regAF, value)
	case F:
		setLowByte(&c.regAF, value)
	case B:
		setHighByte(&c.regBC, value)
	case C:
		setLowByte(&c.regBC, value)
	case D:
		setHighByte(&c.regDE, value)
	case E:
		setLowByte(&c.regDE, value)
	case H:
		setHighByte(&c.regHL, value)
	case L:
		setLowByte(&c.regHL, value)
	default:
		panic(fmt.Sprintf("Invalid register for set byte: %s", reg.String()))
	}
}

func getHighByte(reg uint16) byte {
	return byte(reg >> 8)
}

func getLowByte(reg uint16) byte {
	return byte(reg)
}

func setHighByte(reg *uint16, value byte) {
	*reg = *reg &^ 0xFF00
	var final uint16 = uint16(value) << 8
	*reg = *reg | final
}

func setLowByte(reg *uint16, value byte) {
	*reg = *reg &^ 0x00FF
	*reg = *reg | uint16(value)
}

func (c *CPU) setFlagZ(value bool) { c.setRegBit(F, 7, value) }
func (c *CPU) setFlagN(value bool) { c.setRegBit(F, 6, value) }
func (c *CPU) setFlagH(value bool) { c.setRegBit(F, 5, value) }
func (c *CPU) setFlagC(value bool) { c.setRegBit(F, 4, value) }

func (c *CPU) setRegBit(reg register, bit int, value bool) {
	if reg != F {
		panic(fmt.Sprintf("Unexpected register for get set bit: %s", reg.String()))
	}

	c.setRegByte(reg, memory.SetBit(c.getRegByte(reg), bit, value))
}

func (c *CPU) debugLog() {
	//msg := c.Debug()
	//	fmt.Print(msg)
	//if c.debugLogFile != nil {
	// TODO - handle the PC address being 0xFFFF so trying to read would go past the end
	p1 := c.memory.ReadByte(c.regPC)
	p2 := c.memory.ReadByte(c.regPC + 1)
	p3 := c.memory.ReadByte(c.regPC + 2)
	p4 := c.memory.ReadByte(c.regPC + 3)

	msg := fmt.Sprintf(
		"[A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X]",
		c.getRegByte(A),
		c.getRegByte(F),
		c.getRegByte(B),
		c.getRegByte(C),
		c.getRegByte(D),
		c.getRegByte(E),
		c.getRegByte(H),
		c.getRegByte(L),
		c.getRegShort(SP),
		c.getRegShort(PC),
		p1,
		p2,
		p3,
		p4)

	c.log.Debug(msg)
	//}
}

func (c *CPU) Debug() string {

	// TODO - handle the PC address being 0xFFFF so trying to read would go past the end
	p1 := c.memory.ReadByte(c.regPC)
	p2 := c.memory.ReadByte(c.regPC + 1)
	p3 := c.memory.ReadByte(c.regPC + 2)
	p4 := c.memory.ReadByte(c.regPC + 3)

	name := opcodeNames[p1]
	if p1 == 0xCB {
		// TODO - is p2 right?1
		name = name + " - " + cbNames[p2]
	}

	return fmt.Sprintf(
		"%d=>Next_Opcode:0x%02X %12s -> A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X Z:%t N:%t H:%t C:%t\n",
		c.count,
		p1,
		name,
		c.getRegByte(A),
		c.getRegByte(F),
		c.getRegByte(B),
		c.getRegByte(C),
		c.getRegByte(D),
		c.getRegByte(E),
		c.getRegByte(H),
		c.getRegByte(L),
		c.getRegShort(SP),
		c.getRegShort(PC),
		p1,
		p2,
		p3,
		p4,
		c.getFlagZ(),
		c.getFlagN(),
		c.getFlagH(),
		c.getFlagC())
}

func (c *CPU) DumpOpcodesUsed() {
	f, _ := os.Create("./opcodes_used.txt")
	for x := 0; x < len(c.opcodesUsed); x++ {
		if !c.opcodesUsed[x] {
			continue
		}

		f.WriteString(fmt.Sprintf("0x%02X\n - %s", x, opcodeNames[x]))
	}
	f.Close()
}

func (c *CPU) op_LD_B_B() { c.loadByte(B, B) }
func (c *CPU) op_LD_B_C() { c.loadByte(B, C) }
func (c *CPU) op_LD_B_D() { c.loadByte(B, D) }
func (c *CPU) op_LD_B_E() { c.loadByte(B, E) }
func (c *CPU) op_LD_B_H() { c.loadByte(B, H) }
func (c *CPU) op_LD_B_L() { c.loadByte(B, L) }
func (c *CPU) op_LD_B_A() { c.loadByte(B, A) }
func (c *CPU) op_LD_C_B() { c.loadByte(C, D) }
func (c *CPU) op_LD_C_C() { c.loadByte(C, C) }
func (c *CPU) op_LD_C_D() { c.loadByte(C, D) }
func (c *CPU) op_LD_C_E() { c.loadByte(C, E) }
func (c *CPU) op_LD_C_H() { c.loadByte(C, H) }
func (c *CPU) op_LD_C_L() { c.loadByte(C, L) }
func (c *CPU) op_LD_C_A() { c.loadByte(C, A) }
func (c *CPU) op_LD_D_B() { c.loadByte(D, B) }
func (c *CPU) op_LD_D_C() { c.loadByte(D, C) }
func (c *CPU) op_LD_D_D() { c.loadByte(D, D) }
func (c *CPU) op_LD_D_E() { c.loadByte(D, E) }
func (c *CPU) op_LD_D_H() { c.loadByte(D, H) }
func (c *CPU) op_LD_D_L() { c.loadByte(D, L) }
func (c *CPU) op_LD_D_A() { c.loadByte(D, A) }
func (c *CPU) op_LD_E_B() { c.loadByte(E, B) }
func (c *CPU) op_LD_E_C() { c.loadByte(E, C) }
func (c *CPU) op_LD_E_D() { c.loadByte(E, D) }
func (c *CPU) op_LD_E_E() { c.loadByte(E, E) }
func (c *CPU) op_LD_E_H() { c.loadByte(E, H) }
func (c *CPU) op_LD_E_L() { c.loadByte(E, L) }
func (c *CPU) op_LD_E_A() { c.loadByte(E, A) }
func (c *CPU) op_LD_H_B() { c.loadByte(H, B) }
func (c *CPU) op_LD_H_C() { c.loadByte(H, C) }
func (c *CPU) op_LD_H_D() { c.loadByte(H, D) }
func (c *CPU) op_LD_H_E() { c.loadByte(H, E) }
func (c *CPU) op_LD_H_H() { c.loadByte(H, H) }
func (c *CPU) op_LD_H_L() { c.loadByte(H, L) }
func (c *CPU) op_LD_H_A() { c.loadByte(H, A) }
func (c *CPU) op_LD_L_B() { c.loadByte(L, B) }
func (c *CPU) op_LD_L_C() { c.loadByte(L, C) }
func (c *CPU) op_LD_L_D() { c.loadByte(L, D) }
func (c *CPU) op_LD_L_E() { c.loadByte(L, E) }
func (c *CPU) op_LD_L_H() { c.loadByte(L, H) }
func (c *CPU) op_LD_L_L() { c.loadByte(L, L) }
func (c *CPU) op_LD_L_A() { c.loadByte(L, A) }
func (c *CPU) op_LD_A_B() { c.loadByte(A, B) }
func (c *CPU) op_LD_A_C() { c.loadByte(A, C) }
func (c *CPU) op_LD_A_D() { c.loadByte(A, D) }
func (c *CPU) op_LD_A_E() { c.loadByte(A, E) }
func (c *CPU) op_LD_A_H() { c.loadByte(A, H) }
func (c *CPU) op_LD_A_L() { c.loadByte(A, L) }
func (c *CPU) op_LD_A_A() { c.loadByte(A, A) }

func (c *CPU) loadByte(dest register, src register) {
	c.setRegByte(dest, c.getRegByte(src))
}

func (c *CPU) op_LD_B_n() { c.setRegFromMemory(B) }
func (c *CPU) op_LD_D_n() { c.setRegFromMemory(D) }
func (c *CPU) op_LD_H_n() { c.setRegFromMemory(H) }
func (c *CPU) op_LD_C_n() { c.setRegFromMemory(C) }
func (c *CPU) op_LD_E_n() { c.setRegFromMemory(E) }
func (c *CPU) op_LD_L_n() { c.setRegFromMemory(L) }
func (c *CPU) op_LD_A_n() { c.setRegFromMemory(A) }

func (c *CPU) setRegFromMemory(reg register) {
	c.setRegByte(reg, c.memory.ReadByte(c.regPC))
	c.regPC++
}

func (c *CPU) op_LD_BC_nn() {
	c.setRegShort(BC, c.memory.ReadShort(c.regPC))
	c.regPC += 2
}
func (c *CPU) op_LD_DE_nn() {
	c.setRegShort(DE, c.memory.ReadShort(c.regPC))
	c.regPC += 2
}
func (c *CPU) op_LD_HL_nn() {
	c.setRegShort(HL, c.memory.ReadShort(c.regPC))
	c.regPC += 2
}
func (c *CPU) op_LD_SP_nn() {
	c.setRegShort(SP, c.memory.ReadShort(c.regPC))
	c.regPC += 2
}
func (c *CPU) op_LD_A_HL_plus() {
	hl := c.getRegShort(HL)
	c.setRegByte(A, c.memory.ReadByte(hl))
	hl++
	c.setRegShort(HL, hl)
}
func (c *CPU) op_LD_DE_A() {
	c.memory.WriteByte(c.getRegShort(DE), c.getRegByte(A))
}
func (c *CPU) op_LD_nn_A() {
	c.memory.WriteByte(c.regPC, c.getRegByte(A))
	c.regPC += 2
}
func (c *CPU) op_LDH_n_A() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	var addr uint16 = 0xFF00 | uint16(n)
	c.memory.WriteByte(addr, c.getRegByte(A))
}

func (c *CPU) op_LD_HL_B() { c.LD_HL_x(B) }
func (c *CPU) op_LD_HL_C() { c.LD_HL_x(C) }
func (c *CPU) op_LD_HL_D() { c.LD_HL_x(D) }
func (c *CPU) op_LD_HL_E() { c.LD_HL_x(E) }
func (c *CPU) op_LD_HL_H() { c.LD_HL_x(H) }
func (c *CPU) op_LD_HL_L() { c.LD_HL_x(L) }
func (c *CPU) op_LD_HL_A() { c.LD_HL_x(A) }

func (c *CPU) LD_HL_x(reg register) {
	c.memory.WriteByte(c.regHL, c.getRegByte(reg))
}

func (c *CPU) op_LD_A_DE() {
	c.setRegByte(A, c.memory.ReadByte(c.getRegShort(DE)))
}

func (c *CPU) op_LD_HL_plus_A() {
	c.memory.WriteByte(c.regHL, c.getRegByte(A))
	c.regHL++
}

func (c *CPU) op_LD_HL_sub_A() {
	c.memory.WriteByte(c.regHL, c.getRegByte(A))
	c.regHL--
}

func (c *CPU) op_LDH_C_A() {
	addr := 0xFF00 | uint16(c.getRegByte(C))
	c.memory.WriteByte(addr, c.getRegByte(A))
}

func (c *CPU) op_LD_B_HL() { c.LD_x_HL(B) }
func (c *CPU) op_LD_D_HL() { c.LD_x_HL(D) }
func (c *CPU) op_LD_H_HL() { c.LD_x_HL(H) }
func (c *CPU) op_LD_C_HL() { c.LD_x_HL(C) }
func (c *CPU) op_LD_E_HL() { c.LD_x_HL(E) }
func (c *CPU) op_LD_L_HL() { c.LD_x_HL(L) }
func (c *CPU) op_LD_A_HL() { c.LD_x_HL(A) }

func (c *CPU) LD_x_HL(reg register) {
	c.setRegByte(reg, c.memory.ReadByte(c.regHL))
}

func (c *CPU) op_INC_C()  { c.incrementByteRegister(C) }
func (c *CPU) op_INC_E()  { c.incrementByteRegister(E) }
func (c *CPU) op_INC_L()  { c.incrementByteRegister(L) }
func (c *CPU) op_INC_A()  { c.incrementByteRegister(A) }
func (c *CPU) op_INC_B()  { c.incrementByteRegister(B) }
func (c *CPU) op_INC_D()  { c.incrementByteRegister(D) }
func (c *CPU) op_INC_H()  { c.incrementByteRegister(H) }
func (c *CPU) op_DEC_C()  { c.decrementByteRegister(C) }
func (c *CPU) op_DEC_E()  { c.decrementByteRegister(E) }
func (c *CPU) op_DEC_L()  { c.decrementByteRegister(L) }
func (c *CPU) op_DEC_A()  { c.decrementByteRegister(A) }
func (c *CPU) op_DEC_B()  { c.decrementByteRegister(B) }
func (c *CPU) op_DEC_D()  { c.decrementByteRegister(D) }
func (c *CPU) op_DEC_H()  { c.decrementByteRegister(H) }
func (c *CPU) op_INC_BC() { c.incrementShortRegister(BC) }
func (c *CPU) op_INC_DE() { c.incrementShortRegister(DE) }
func (c *CPU) op_INC_HL() { c.incrementShortRegister(HL) }
func (c *CPU) op_INC_SP() { c.incrementShortRegister(SP) }
func (c *CPU) op_DEC_BC() { c.decrementShortRegister(BC) }
func (c *CPU) op_DEC_DE() { c.decrementShortRegister(DE) }
func (c *CPU) op_DEC_HL() { c.decrementShortRegister(HL) }
func (c *CPU) op_DEC_SP() { c.decrementShortRegister(SP) }

func (c *CPU) incrementByteRegister(reg register) {
	value := c.getRegByte(reg)
	current := value + 1
	c.setFlagZ(current == 0)
	c.setFlagN(false)

	abc := current ^ 0x01 ^ value
	c.setFlagH(abc&0x10 == 0x10) //((current - 1) & 0x10) == 0x10)
	//c.setFlagC(true) // TODO - why?
	c.setRegByte(reg, current)
}

func (c *CPU) decrementByteRegister(reg register) {
	value := c.getRegByte(reg)
	current := value - 1
	c.setFlagZ(current == 0)
	c.setFlagN(true)

	abc := current ^ 0x01 ^ value
	c.setFlagH(abc&0x10 == 0x10) //((current + 1) & 0x10) == 0x10)
	c.setRegByte(reg, current)
}

func (c *CPU) incrementShortRegister(reg register) {
	current := c.getRegShort(reg) + 1
	//	c.setFlagZ(current == 0)
	//	c.setFlagN(true) // TODO - why?
	//	c.setFlagH(current&0x10 == 0x10)
	//	c.setFlagC(true) // TODO - why?
	c.setRegShort(reg, current)
}

func (c *CPU) decrementShortRegister(reg register) {
	current := c.getRegShort(reg) + 1
	//	c.setFlagZ(current == 0)
	//	c.setFlagN(true)
	//	c.setFlagH(current&0x10 == 0x10)
	c.setRegShort(reg, current)
}

func (c *CPU) op_ADD_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.op_ADD(n)
}

func (c *CPU) op_ADD(n byte) {
	a := c.getRegByte(A)
	result := a + n
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH((a & 0x10) == 0x10)
	c.setFlagC((a & 0x0F) == 0x0F)
}

func (c *CPU) op_SUB_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.op_SUB(n)
}

func (c *CPU) op_SUB(n byte) {
	a := c.getRegByte(A)
	result := a - n
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(true)
	c.setFlagH((a & 0x10) == 0x10)
	c.setFlagC((a & 0x0F) == 0x0F)
}

func (c *CPU) op_NOP() {}

func (c *CPU) op_DI() {
	// TODO - set IME?
}

func (c *CPU) op_JP() {
	newAddr := c.memory.ReadShort(c.regPC)
	c.regPC = newAddr
}

func (c *CPU) op_JR_e() {
	c.conditionalJumpToOffset(true)
}

func (c *CPU) op_JR_NZ_e() {
	c.conditionalJumpToOffset(!c.getFlagZ())
}

func (c *CPU) op_JR_NC_e() {
	c.conditionalJumpToOffset(!c.getFlagC())
}

func (c *CPU) op_JR_Z_e() {
	c.conditionalJumpToOffset(c.getFlagZ())
}

func (c *CPU) op_JR_C_e() {
	c.conditionalJumpToOffset(c.getFlagC())
}

func (c *CPU) conditionalJumpToOffset(condition bool) {
	offset := c.memory.ReadByte((c.regPC))
	c.regPC++

	if !condition {
		return
	}

	//	fmt.Printf("JMP PC: 0x%04X, offset: %d, signed: %d", c.regPC, offset, int8(offset))

	// TODO - is this right?
	c.regPC = uint16(int16(c.regPC) + int16(int8(offset)))
}

func (c *CPU) op_CALL_nn() {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	c.regSP -= 2
	c.memory.WriteShort(c.regSP, c.regPC)
	c.regPC = nn
}

func (c *CPU) op_CALL_NZ_nn() { c.callCondition(!c.getFlagZ()) }
func (c *CPU) op_CALL_Z_nn()  { c.callCondition(c.getFlagZ()) }
func (c *CPU) op_CALL_NC_nn() { c.callCondition(!c.getFlagC()) }
func (c *CPU) op_CALL_C_nn()  { c.callCondition(c.getFlagC()) }

func (c *CPU) callCondition(condition bool) {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	if condition {
		c.regSP -= 2
		c.memory.WriteShort(c.regSP, c.regPC)
		c.setRegShort(PC, nn)
	}
}

func (c *CPU) op_RET() {

	c.setRegShort(PC, c.memory.ReadShort(c.getRegShort(SP)))
	c.regSP += 2
}

func (c *CPU) op_PUSH_BC() { c.push(BC) }
func (c *CPU) op_PUSH_DE() { c.push(DE) }
func (c *CPU) op_PUSH_HL() { c.push(HL) }
func (c *CPU) op_PUSH_AF() { c.push(AF) }
func (c *CPU) op_POP_BC()  { c.pop(BC) }
func (c *CPU) op_POP_DE()  { c.pop(DE) }
func (c *CPU) op_POP_HL()  { c.pop(HL) }
func (c *CPU) op_POP_AF()  { c.pop(AF) }

func (c *CPU) push(reg register) {
	c.regSP -= 2
	c.memory.WriteShort(c.regSP, c.getRegShort(reg))
}

func (c *CPU) pop(reg register) {
	c.setRegShort(reg, c.memory.ReadShort(c.regSP))
	c.regSP += 2
}

func (c *CPU) op_OR_B() { c.or(B) }
func (c *CPU) op_OR_C() { c.or(C) }
func (c *CPU) op_OR_D() { c.or(D) }
func (c *CPU) op_OR_E() { c.or(E) }
func (c *CPU) op_OR_H() { c.or(H) }
func (c *CPU) op_OR_L() { c.or(L) }
func (c *CPU) op_OR_A() { c.or(A) }

func (c *CPU) or(reg register) {
	result := c.getRegByte(A) | c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_AND_B() { c.and(B) }
func (c *CPU) op_AND_C() { c.and(C) }
func (c *CPU) op_AND_D() { c.and(D) }
func (c *CPU) op_AND_E() { c.and(E) }
func (c *CPU) op_AND_H() { c.and(H) }
func (c *CPU) op_AND_L() { c.and(L) }
func (c *CPU) op_AND_A() { c.and(A) }

func (c *CPU) and(reg register) {
	result := c.getRegByte(A) & c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_XOR_B() { c.xor(B) }
func (c *CPU) op_XOR_C() { c.xor(C) }
func (c *CPU) op_XOR_D() { c.xor(D) }
func (c *CPU) op_XOR_E() { c.xor(E) }
func (c *CPU) op_XOR_H() { c.xor(H) }
func (c *CPU) op_XOR_L() { c.xor(L) }
func (c *CPU) op_XOR_A() { c.xor(A) }

func (c *CPU) xor(reg register) {
	result := c.getRegByte(A) ^ c.getRegByte(reg)
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_XOR_HL() {
	data := c.memory.ReadByte(c.getRegShort(HL))
	result := c.getRegByte(A) ^ data
	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_CP_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.op_CP(n)
}

func (c *CPU) op_CP(n byte) {
	a := c.getRegByte(A)
	result := a - n
	c.setFlagZ(result == 0)
	c.setFlagN(true)
	c.setFlagH((a & 0x0F) < (n & 0x0F))
	c.setFlagC(a < n)
}

func (c *CPU) op_LDH_A_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	c.setRegByte(A, c.memory.ReadByte(0xFF00|uint16(n)))
}

func (c *CPU) op_LD_A_nn() {
	nn := c.memory.ReadShort(c.regPC)
	c.regPC += 2
	c.setRegByte(A, c.memory.ReadByte(nn))
}

func (c *CPU) op_AND_n() {
	n := c.memory.ReadByte(c.regPC)
	c.regPC++
	result := c.getRegByte(A) & n
	c.setRegByte(A, n)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(true)
	c.setFlagC(false)
}

func (c *CPU) op_CB_op() {
	op := c.memory.ReadByte(c.regPC)
	c.regPC++

	if op >= 0x40 && op <= 0x7F {
		c.op_BIT(op)
	} else if op >= 0x30 && op <= 0x37 {
		c.op_SWAP(op)
	} else if op >= 0x10 && op <= 0x17 {
		c.op_RL(op)
	} else {
		panic(fmt.Sprintf("Unknown op for CB: 0x%02X", op))
	}
}

func (c *CPU) op_BIT(op byte) {
	// TODO - no idea if this is right
	var reg register
	switch op {
	case 0x40:
	case 0x48:
	case 0x50:
	case 0x58:
	case 0x60:
	case 0x68:
	case 0x70:
	case 0x78:
		reg = B
	case 0x41:
	case 0x49:
	case 0x51:
	case 0x59:
	case 0x61:
	case 0x69:
	case 0x71:
	case 0x79:
		reg = C
	case 0x42:
	case 0x4A:
	case 0x52:
	case 0x5A:
	case 0x62:
	case 0x6A:
	case 0x72:
	case 0x7A:
		reg = D
	case 0x43:
	case 0x4B:
	case 0x53:
	case 0x5B:
	case 0x63:
	case 0x6B:
	case 0x73:
	case 0x7B:
		reg = E
	case 0x44:
	case 0x4C:
	case 0x54:
	case 0x5C:
	case 0x64:
	case 0x6C:
	case 0x74:
	case 0x7C:
		reg = H
	case 0x45:
	case 0x4D:
	case 0x55:
	case 0x5D:
	case 0x65:
	case 0x6D:
	case 0x75:
	case 0x7D:
		reg = L
	case 0x46:
	case 0x4E:
	case 0x56:
	case 0x5E:
	case 0x66:
	case 0x6E:
	case 0x76:
	case 0x7E:
		reg = HL
	case 0x47:
	case 0x4F:
	case 0x57:
	case 0x5F:
	case 0x67:
	case 0x6F:
	case 0x77:
	case 0x7F:
		reg = A
	default:
		panic(fmt.Sprintf("Unknown register op for BIT: 0x%02X", op))
		// TODO - implement loads of things
		return
	}

	var bit int
	switch {
	case op >= 0x40 && op <= 0x47:
		bit = 0
	case op >= 0x48 && op <= 0x4F:
		bit = 1
	case op >= 0x50 && op <= 0x57:
		bit = 2
	case op >= 0x58 && op <= 0x5F:
		bit = 3
	case op >= 0x60 && op <= 0x67:
		bit = 4
	case op >= 0x68 && op <= 0x6F:
		bit = 5
	case op >= 0x70 && op <= 0x77:
		bit = 6
	case op >= 0x78 && op <= 0x7F:
		bit = 7
	default:
		panic(fmt.Sprintf("Unknown bit op for BIT: 0x%02X", op))
	}

	value := c.getRegBit(reg, bit)
	// TODO - set Z if false. Not clear if that means set do nothing if
	// value is true though
	c.setFlagZ(!value)
	c.setFlagN(false)
	c.setFlagH(true)
}

func (c *CPU) op_SWAP(op byte) {
	// TODO - no idea if this is right
	var reg register
	switch op {
	case 0x37:
		reg = A
	case 0x30:
		reg = B
	case 0x31:
		reg = C
	case 0x32:
		reg = D
	case 0x33:
		reg = E
	case 0x34:
		reg = H
	case 0x35:
		reg = L
	case 0x36:
		reg = HL
	default:
		panic(fmt.Sprintf("Unknown reg op for SWAP: 0x%02X", op))
	}

	n := c.getRegByte(reg)
	// TODO - swap upper and lower nibbles
	c.setRegByte(reg, n)

	c.setFlagZ(n == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(false)
}

func (c *CPU) op_RL(op byte) {
	var reg register
	switch op {
	case 0x10:
		reg = B
	case 0x11:
		reg = C
	case 0x12:
		reg = D
	case 0x13:
		reg = E
	case 0x14:
		reg = H
	case 0x15:
		reg = L
	case 0x16:
		reg = HL
	case 0x17:
		reg = A
	default:
		panic(fmt.Sprintf("Unknown reg op for RL: 0x%02X", op))
	}

	//oldBit7 := c.getRegBit(reg, 7)

	c.rotateLeft(reg)
	//value := c.getRegByte(reg)

	//c.setFlagZ(value == 0)
	//c.setFlagN(false)
	//c.setFlagH(false)
	//c.setFlagC(oldBit7)
}

func (c *CPU) op_RL_A() {
	result := c.getRegByte(A)

	carry := result & 0x80
	//carry = carry >> 7
	result = result << 1
	//result = result & 0xFE
	//result = result | carry

	if c.getFlagC() {
		result = result ^ 0x01
	}

	c.setRegByte(A, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(carry == 0x80) //carry&0x01 == 0x01)
}

func (c *CPU) rotateLeft(reg register) {
	result := c.getRegByte(reg)

	carry := result & 0x80
	//carry = carry >> 7
	result = result << 1
	//result = result & 0xFE
	//result = result | carry
	if c.getFlagC() {
		result = result ^ 0x01
	}

	c.setRegByte(reg, result)
	c.setFlagZ(result == 0)
	c.setFlagN(false)
	c.setFlagH(false)
	c.setFlagC(carry == 0x80) //carry&0x01 == 0x01)
}

func (c *CPU) op_SUB_B()  { c.op_SUB(c.getRegByte(B)) }
func (c *CPU) op_SUB_C()  { c.op_SUB(c.getRegByte(C)) }
func (c *CPU) op_SUB_D()  { c.op_SUB(c.getRegByte(D)) }
func (c *CPU) op_SUB_E()  { c.op_SUB(c.getRegByte(E)) }
func (c *CPU) op_SUB_H()  { c.op_SUB(c.getRegByte(H)) }
func (c *CPU) op_SUB_L()  { c.op_SUB(c.getRegByte(L)) }
func (c *CPU) op_SUB_HL() { c.op_SUB(c.memory.ReadByte(c.getRegShort(HL))) }
func (c *CPU) op_SUB_A()  { c.op_SUB(c.getRegByte(A)) }

func (c *CPU) op_CP_B()  { c.op_CP(c.getRegByte(B)) }
func (c *CPU) op_CP_C()  { c.op_CP(c.getRegByte(C)) }
func (c *CPU) op_CP_D()  { c.op_CP(c.getRegByte(D)) }
func (c *CPU) op_CP_E()  { c.op_CP(c.getRegByte(E)) }
func (c *CPU) op_CP_H()  { c.op_CP(c.getRegByte(H)) }
func (c *CPU) op_CP_L()  { c.op_CP(c.getRegByte(L)) }
func (c *CPU) op_CP_HL() { c.op_CP(c.memory.ReadByte(c.getRegShort(HL))) }
func (c *CPU) op_CP_A()  { c.op_CP(c.getRegByte(A)) }

func (c *CPU) op_ADD_B()  { c.op_ADD(c.getRegByte(B)) }
func (c *CPU) op_ADD_C()  { c.op_ADD(c.getRegByte(C)) }
func (c *CPU) op_ADD_D()  { c.op_ADD(c.getRegByte(D)) }
func (c *CPU) op_ADD_E()  { c.op_ADD(c.getRegByte(E)) }
func (c *CPU) op_ADD_H()  { c.op_ADD(c.getRegByte(H)) }
func (c *CPU) op_ADD_L()  { c.op_ADD(c.getRegByte(L)) }
func (c *CPU) op_ADD_HL() { c.op_ADD(c.memory.ReadByte(c.getRegShort(HL))) }
func (c *CPU) op_ADD_A()  { c.op_ADD(c.getRegByte(A)) }

func (c *CPU) op_LD_HL_n() {
	pc := c.getRegShort(PC)
	n := c.memory.ReadByte(pc)
	c.setRegShort(PC, pc+1)
	c.memory.WriteByte(c.getRegShort(HL), n)
}
