package intcode

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	positionMode = iota
	intermediateMode
	relativeMode
)
const supportedParamCount = 3

// IntCode defines a IntCode computer
type IntCode struct {
	outputs      []int
	pointer      int
	memory       []int
	paused       bool
	relativeBase int // Defaults to zero

	Print            bool // Indicator if the output should be printed or not
	InputHandler     func() int
	OnOutputRecieved func(output int)
}

// NewIntCode creates a new intcode computer
func NewIntCode() *IntCode {
	return &IntCode{
		InputHandler: DefaultInputHandler,
		outputs:      make([]int, 0),
		Print:        false,
		pointer:      0,
	}
}

// SetMemory set the memory of the machine
func (c *IntCode) SetMemory(memory []int) {
	c.memory = memory
}

// Pause pauses the program
func (c *IntCode) Pause() {
	c.paused = true
}

// Execute executes the program
func (c *IntCode) Execute() (finished bool) {
	c.paused = false

	for !c.paused {

		i := c.pointer

		digits := getDigits(c.memory[i])
		digits = append(digits, make([]int, (2+supportedParamCount)-len(digits))...) // Add unprovided modes as 0s
		opcode := digits[0] + digits[1]*10
		modes := digits[2:]

		switch opcode {
		case 1:
			{
				c.write(c.getWriteParam(i+3, modes[2]), c.getParam(i+1, modes[0])+c.getParam(i+2, modes[1]))
				c.pointer = i + 4
			}
		case 2:
			{
				c.write(c.getWriteParam(i+3, modes[2]), c.getParam(i+1, modes[0])*c.getParam(i+2, modes[1]))
				c.pointer = i + 4
			}
		case 3:
			{
				input := c.InputHandler()
				// The program might have been paused in the input handler
				if c.paused {
					return
				}

				c.write(c.getWriteParam(i+1, modes[0]), input)
				c.pointer = i + 2
			}
		case 4:
			{
				value := c.getParam(i+1, modes[0])
				c.outputs = append(c.outputs, value)
				if c.Print {
					fmt.Println(value)
				}
				c.pointer = i + 2

				if c.OnOutputRecieved != nil {
					c.OnOutputRecieved(value)
				}
			}
		case 5:
			{
				pointer := i + 3
				if c.getParam(i+1, modes[0]) != 0 {
					pointer = c.getParam(i+2, modes[1])
				}
				c.pointer = pointer
			}
		case 6:
			{
				pointer := i + 3
				if c.getParam(i+1, modes[0]) == 0 {
					pointer = c.getParam(i+2, modes[1])
				}
				c.pointer = pointer
			}
		case 7:
			{
				value := 0
				if c.getParam(i+1, modes[0]) < c.getParam(i+2, modes[1]) {
					value = 1
				}
				c.write(c.getWriteParam(i+3, modes[2]), value)
				c.pointer = i + 4
			}
		case 8:
			{
				value := 0
				if c.getParam(i+1, modes[0]) == c.getParam(i+2, modes[1]) {
					value = 1
				}
				c.write(c.getWriteParam(i+3, modes[2]), value)
				c.pointer = i + 4
			}
		case 9:
			{
				value := c.getParam(i+1, modes[0])
				c.relativeBase += value
				c.pointer = i + 2
			}
		case 99:
			{
				return true // The program finished
			}
		default:
			{
				log.Fatalf("Invalid opcode: %d at %d\n%+v", opcode, i, c.memory)
			}
		}
	}

	return false // The program did not finish
}

// Output returns the output of the machine
func (c *IntCode) Output() []int {
	return c.outputs
}

// ResetPointer resets the instruction pointer to 0
func (c *IntCode) ResetPointer() {
	c.pointer = 0
}

func (c *IntCode) getParam(pos int, mode int) int {

	var position int
	if mode == positionMode {
		position = c.memory[pos]
	} else if mode == relativeMode {
		position = c.relativeBase + c.memory[pos]
	} else if mode == intermediateMode {
		return c.memory[pos]
	} else {
		log.Fatalf("Invalid mode %d", mode)
	}

	// Check if this parameter will go outside memory
	if position >= len(c.memory) {
		// Extend memory
		c.extendMemory(position - len(c.memory) + 1)
	}

	return c.memory[position]
}

func (c *IntCode) getWriteParam(pos int, mode int) int {
	if mode == positionMode || mode == intermediateMode {
		return c.memory[pos]
	} else if mode == relativeMode {
		return c.relativeBase + c.memory[pos]
	} else {
		log.Fatalf("Invalid mode %d", mode)
	}

	return -1
}

func (c *IntCode) write(pos int, value int) {
	if pos >= len(c.memory) {
		c.extendMemory(pos - len(c.memory) + 1)
	}

	c.memory[pos] = value
}

func (c *IntCode) extendMemory(amount int) {
	// Extend memory
	c.memory = append(c.memory, make([]int, amount)...)
}

func getDigits(value int) []int {
	digits := make([]int, 0)
	for value > 0 {
		digits = append(digits, value%10)
		value /= 10
	}
	return digits
}

func DefaultInputHandler() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter input number: ")
	numberText, err := reader.ReadString('\n')
	input, err := strconv.Atoi(strings.Trim(numberText, "\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	return input
}

// ParseFromFile reads and parses an intcode program from file
func ParseFromFile(fileName string) []int {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	strValues := strings.Split(string(bytes), ",")
	values := make([]int, 0)
	for _, v := range strValues {
		i, err := strconv.Atoi(strings.Trim(v, "\n"))
		if err != nil {
			log.Fatal(err)
		}
		values = append(values, i)
	}

	return values
}
