/*
What I am trying to to do
  - Creating a REPL which takes input from the CLI
    (READ EVALUATE PRINT LOOP )
  - Then We are going to tokenize the input from the user after each word
    (meaning: After each space we consider it as a token)
*/
package REPL

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// initialze the capacity of the buffer to 1KB
var defaultBufferCapacity uint16 = 1024

// Stores the input of the user into the buffer
type LineBuffer struct {
	Buffer         []byte // Where we store the data
	bufferCapacity uint16 // The capacity of the buffer we are trying to store
	BufferLength   uint16
	pos            uint16 //gives the information about the position of the buffer
}

// Just prints godb in each line
func PrintDB() {
	fmt.Printf("godb >")
}

// Initializing the buffer  f1
func InitLineBuffer() *LineBuffer {
	return &LineBuffer{
		Buffer:         make([]byte, 0, defaultBufferCapacity),
		bufferCapacity: defaultBufferCapacity,
		BufferLength:   0,
	}
}

// Taking input from the user f2
func (lb *LineBuffer) UserInput() {
	reader := bufio.NewReader(os.Stdin) // a struct

	// Tells the user to take input in the buffer
	// until the occurence of first newline
	input, _ := reader.ReadBytes(';')
	if len(input) > 0 && input[len(input)-1] == ';' { // checking if the input is not empty and ends with a semicolon
		input = input[:len(input)] //assigning the input to the buffer
	}
	input = bytes.TrimRight(input, "\n")
	lb.Write(input)
}

// Writing the input from the user into the buffer
func (lb *LineBuffer) Write(input []byte) {
	// input... is a variadic expression- it unpacks the input slice into
	// indivisual bytes that can be appended
	lb.Buffer = append(lb.Buffer, input...)
	lb.BufferLength = uint16(len(lb.Buffer))
}

func (lb *LineBuffer) Output() {
	fmt.Printf("%s\n", lb.Buffer)
}

// Resets the buffer
func (lb *LineBuffer) Reset() {
	lb.Buffer = lb.Buffer[:0] // keeps the alloted memory but clears contents
	lb.BufferLength = 0
	lb.pos = 0
}

/*
	The flow of the program

		lineBuffer := InitLineBuffer() // initializes file buffer
		PrintDB() // prints the db > on the screen
		lineBuffer.UserInput() // takes input from the User form stdin(keyboard)
		lineBuffer.Output()	 // for now just prints what is in the buffer
		lineBuffer.Reset() // resets the buffer or empties the buffer
*/
