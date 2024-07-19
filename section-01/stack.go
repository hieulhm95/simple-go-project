package main

type Stack struct {
	data []interface{}
}

// Create a new stack
func New() *Stack {
	// TODO
	return &Stack{
		data: make([]interface{}, 0),
	}
}

// Return the number of items in the stack
func (this *Stack) Len() int {
	// TODO
	return len(this.data)
}

// View the top item on the stack
func (this *Stack) Peek() interface{} {
	// TODO
	if len(this.data) > 0 {
		return this.data[len(this.data)-1]
	}
	return nil

}

// Pop the top item of the stack and return it
// func (this *Stack) Pop() interface{} {
// 	// TODO
// }

// Push a value onto the top of the stack
func (this *Stack) Push(value interface{}) {
	// TODO
	this.data = append(this.data, value)

}
