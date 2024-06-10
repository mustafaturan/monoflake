package monoflake_test

import (
	"fmt"

	"github.com/mustafaturan/monoflake"
)

func ExampleMonoFlake_Next() {
	var nodeID uint16 = 19
	tid, err := monoflake.New(nodeID)

	if err != nil {
		panic(err)
	}
	i1, i2 := tid.Next(), tid.Next()
	fmt.Println(i1 < i2)
	// Output:
	// true
}

func ExampleMonoFlake_NextBase62() {
	var nodeID uint16 = 19
	tid, err := monoflake.New(nodeID)

	if err != nil {
		panic(err)
	}
	fmt.Println(len(tid.Next().String()))
	// Output:
	// 11
}

func ExampleMonoFlake_NextBytes() {
	var nodeID uint16 = 19
	tid, err := monoflake.New(nodeID)

	if err != nil {
		panic(err)
	}
	fmt.Println(len(tid.Next().Bytes()))
	// Output:
	// 11
}
