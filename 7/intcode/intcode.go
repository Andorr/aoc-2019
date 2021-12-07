package intcode

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	positionMode = iota
	intermediateMode
)
const supportedParamCount = 3

// IntCode defines a IntCode computer
type IntCode struct {
	outputs []int
	pointer int
	memory  []int
	paused  bool

	Print            bool // Indicator if the output should be printed or not
	InputHandler     func() int
	OnOutputRecieved func(output int)
}

// NewIntCode creates a new intcode computer
func NewIntCode() *IntCode {
	return &IntCode{
		InputHandler: inputHandler,
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
	program := c.memory

	for !c.paused {

		i := c.pointer

		// Get the current program
		digits := getDigits(program[i])
		digits = append(digits, make([]int, (2+supportedParamCount)-len(digits))...) // Add unprovided modes as 0s
		opcode := digits[0] + digits[1]*10
		modes := digits[2:]

		switch opcode {
		case 1:
			{
				program[program[i+3]] = getParam(program, i+1, modes[0]) + getParam(program, i+2, modes[1])
				c.pointer = i + 4
			}
		case 2:
			{
				program[program[i+3]] = getParam(program, i+1, modes[0]) * getParam(program, i+2, modes[1])
				c.pointer = i + 4
			}
		case 3:
			{
				input := c.InputHandler()
				// The program might have been paused in the input handler
				if c.paused {
					return
				}
				program[getParam(program, i+1, 1)] = input
				c.pointer = i + 2
			}
		case 4:
			{
				value := getParam(program, i+1, modes[0])
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
				if getParam(program, i+1, modes[0]) != 0 {
					pointer = getParam(program, i+2, modes[1])
				}
				c.pointer = pointer
			}
		case 6:
			{
				pointer := i + 3
				if getParam(program, i+1, modes[0]) == 0 {
					pointer = getParam(program, i+2, modes[1])
				}
				c.pointer = pointer
			}
		case 7:
			{
				value := 0
				if getParam(program, i+1, modes[0]) < getParam(program, i+2, modes[1]) {
					value = 1
				}
				program[getParam(program, i+3, 1)] = value
				c.pointer = i + 4
			}
		case 8:
			{
				value := 0
				if getParam(program, i+1, modes[0]) == getParam(program, i+2, modes[1]) {
					value = 1
				}
				program[getParam(program, i+3, 1)] = value
				c.pointer = i + 4
			}
		case 99:
			{
				return true // The program finished
			}
		default:
			{
				log.Fatalf("Invalid opcode: %d at %d\n%+v", opcode, i, program)
			}
		}
	}

	return false // The program did not finish
}

// Outputs returns the output of the machine
func (c *IntCode) Outputs() []int {
	return c.outputs
}

// ResetPointer resets the instruction pointer to 0
func (c *IntCode) ResetPointer() {
	c.pointer = 0
}

func getParam(program []int, pos int, mode int) int {

	if mode == positionMode {
		return program[program[pos]]
	} else if mode == intermediateMode {
		return program[pos]
	}
	log.Fatalf("Invalid mode %d", mode)
	return -1
}

func getDigits(value int) []int {
	digits := make([]int, 0)
	for value > 0 {
		digits = append(digits, value%10)
		value /= 10
	}
	return digits
}

func inputHandler() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter input number: ")
	numberText, err := reader.ReadString('\n')
	input, err := strconv.Atoi(strings.Trim(numberText, "\r\n"))
	if err != nil {
		log.Fatal(err)
	}
	return input
}
