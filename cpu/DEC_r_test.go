package cpu

import (
	"testing"

	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
	"github.com/stretchr/testify/suite"
)

func Test_DEC_B(t *testing.T) {
	suite.Run(t, &decTestSuite{reg: B})
}

func Test_DEC_D(t *testing.T) {
	suite.Run(t, &decTestSuite{reg: D})
}

func Test_DEC_H(t *testing.T) {
	suite.Run(t, &decTestSuite{reg: H})
}

func Test_DEC_C(t *testing.T) {
	suite.Run(t, &decTestSuite{reg: C})
}

func Test_DEC_E(t *testing.T) {
	suite.Run(t, &decTestSuite{reg: E})
}

func Test_DEC_L(t *testing.T) {
	suite.Run(t, &decTestSuite{reg: L})
}

func Test_DEC_A(t *testing.T) {
	suite.Run(t, &decTestSuite{reg: A})
}

func Benchmark_DEC(b *testing.B) {
	opcode := createDEC_r(0x00, B)
	regs := &Registers{}
	mem := memory.CreateMemory(&log.Log{})

	for x := 0; x < b.N; x++ {
		opcode.doCycle(1, regs, mem)
	}
}

type decTestSuite struct {
	suite.Suite
	reg Register
}

func (i *decTestSuite) test(initial uint8, expected uint8, carry bool) {
	opcode := createDEC_r(0x00, i.reg)

	regs := &testRegisters_UseOneRegister{
		test:         i.Suite.T(),
		allowedReg:   i.reg,
		allowedFlags: []RegisterFlags{ZFlag, NFlag, HFlag},
	}
	mem := &testMemory_NoAccess{test: i.Suite.T()}

	regs.Set8(i.reg, initial)

	completed, err := opcode.doCycle(1, regs, mem)

	i.Nil(err)
	i.Equal(expected, regs.Get8(i.reg))
	i.Equal(carry, regs.GetFlag(HFlag))
	i.Equal(expected == 0, regs.GetFlag(ZFlag))
	i.False(regs.GetFlag(NFlag))

	i.True(completed)
}

func (i *decTestSuite) Test_0() {
	i.test(0, 255, false)
}

func (i *decTestSuite) Test_1() {
	i.test(1, 0, false)
}

func (i *decTestSuite) Test_254() {
	i.test(254, 253, false)
}

func (i *decTestSuite) Test_255() {
	i.test(255, 254, false)
}

func (i *decTestSuite) Test_15() {
	i.test(15, 14, false)
}

func (i *decTestSuite) Test_16() {
	i.test(16, 15, false)
}
