package m68000

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_A(t *testing.T) {
	tests := []struct {
		input    []uint8
		expected string
	}{
		{[]uint8{0x00}, "Unknown instruction 0x00"},
	}

	parser := instructionParser{}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			t.Parallel()

			result := parser.parse(test.input, 0)

			if result.Name() != test.expected {
				t.Errorf("Input: 0x%01X, Expected: '%s', Actual: '%s'", test.input, test.expected, result.Name())
			}
		})
	}
}

func Test_Comparison(t *testing.T) {
	const input = ""

	file, err := os.Open("../../../cps1/test/Comparison/" + input + "/input.dmp")

	assert.Nil(t, err, "missing test data file")
	defer file.Close()

	r, _ := regexp.Compile(`\w{4}`)
	code := []uint8{}
	inputScanner := bufio.NewScanner(file)
	for inputScanner.Scan() {
		line := inputScanner.Text()
		matches := r.FindAllString(line, -1)

		for x := 1; x < len(matches); x++ {
			val, _ := strconv.ParseUint(matches[x], 16, 16)
			lower := uint8(val)
			upper := uint8(val >> 8)
			code = append(code, upper)
			code = append(code, lower)
		}
	}

	parser := instructionParser{}

	expected, err := os.ReadFile("../../../cps1/test/Comparison/" + input + "/output.asm")
	assert.Nil(t, err, "Missing expected result file")

	expectedLines := strings.Split(string(expected), "\n")
	i := 0

	for x := 0; x < len(expectedLines); x++ {
		expectedLine := expectedLines[x]

		var resultLine string
		inst := parser.parse(code, i)
		switch inst.Length() {
		case 2:
			resultLine = fmt.Sprintf("%06X: %02X%02X                     %s", i, code[i], code[i+1], inst.Name())
		case 4:
			resultLine = fmt.Sprintf("%06X: %02X%02X %02X%02X                %s", i, code[i], code[i+1], code[i+2], code[i+3], inst.Name())
		case 6:
			resultLine = fmt.Sprintf("%06X: %02X%02X %02X%02X %02X%02X           %s", i, code[i], code[i+1], code[i+2], code[i+3], code[i+4], code[i+5], inst.Name())
		case 8:
			resultLine = fmt.Sprintf("%06X: %02X%02X %02X%02X %02X%02X %02X%02X      %s", i, code[i], code[i+1], code[i+2], code[i+3], code[i+4], code[i+5], code[i+6], code[i+7], inst.Name())
		default:
			panic(fmt.Sprintf("Invalid instruction length of %d", inst.Length()))
		}

		require.Equal(t, expectedLine, resultLine, fmt.Sprintf("Instruction: %08b%08b", code[i], code[i+1]))

		i += inst.Length()
	}
}
