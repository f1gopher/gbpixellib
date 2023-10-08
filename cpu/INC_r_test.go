package cpu

import (
	"testing"

	"github.com/f1gopher/gbpixellib/log"
	"github.com/f1gopher/gbpixellib/memory"
	"github.com/stretchr/testify/suite"
)

func Test_INC_B(t *testing.T) {
	suite.Run(t, &incTestSuite{reg: B})
}

func Test_INC_D(t *testing.T) {
	suite.Run(t, &incTestSuite{reg: D})
}

func Test_INC_H(t *testing.T) {
	suite.Run(t, &incTestSuite{reg: H})
}

func Test_INC_C(t *testing.T) {
	suite.Run(t, &incTestSuite{reg: C})
}

func Test_INC_E(t *testing.T) {
	suite.Run(t, &incTestSuite{reg: E})
}

func Test_INC_L(t *testing.T) {
	suite.Run(t, &incTestSuite{reg: L})
}

func Test_INC_A(t *testing.T) {
	suite.Run(t, &incTestSuite{reg: A})
}

func Benchmark_INC(b *testing.B) {
	opcode := createINC_r(0x00, B)
	regs := &Registers{}
	mem := memory.CreateMemory(&log.Log{})

	for x := 0; x < b.N; x++ {
		opcode.doCycle(1, regs, mem)
	}
}

type incTestSuite struct {
	suite.Suite
	reg register
}

func (i *incTestSuite) test(initial uint8, expected uint8, carry bool) {
	opcode := createINC_r(0x00, i.reg)

	regs := &testRegisters_UseOneRegister{
		test:         i.Suite.T(),
		allowedReg:   i.reg,
		allowedFlags: []registerFlags{ZFlag, NFlag, HFlag},
	}
	mem := &testMemory_NoAccess{test: i.Suite.T()}

	regs.set8(i.reg, initial)

	completed, err := opcode.doCycle(1, regs, mem)

	i.Nil(err)
	i.Equal(expected, regs.Get8(i.reg))
	i.Equal(carry, regs.GetFlag(HFlag))
	i.Equal(expected == 0, regs.GetFlag(ZFlag))
	i.False(regs.GetFlag(NFlag))

	i.True(completed)
}

func (i *incTestSuite) Test_0() {
	i.test(0, 1, false)
}

func (i *incTestSuite) Test_1() {
	i.test(1, 2, false)
}

func (i *incTestSuite) Test_254() {
	i.test(254, 255, false)
}

//func (i *incTestSuite) Test_255() {
//	i.test(255, 0, false)
//}

func (i *incTestSuite) Test_14() {
	i.test(14, 15, false)
}

func (i *incTestSuite) Test_15() {
	i.test(15, 16, true)
}
