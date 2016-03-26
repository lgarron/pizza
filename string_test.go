package pizza

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	var p Pizza
	var str string

	// We place known-reported calls at the start and end so we can tell
	// `go vet` is actually looking for them.
	str = fmt.Sprintf("This is used.")
	fmt.Sprintf("This is unused.")

	str = p.makeStringMethod()
	p.makeStringMethod()

	str = p.PublicMakeStringMethod()
	p.PublicMakeStringMethod()

	str = makeStringBareFunction()
	makeStringBareFunction()

	str = PublicMakeStringBareFunction()
	PublicMakeStringBareFunction()

	str = fmt.Sprintf("This is used.")
	fmt.Sprintf("This is unused.")

	// Also try a "reportable" function.
	_, _ = fmt.Printf("This is used.\n")
	fmt.Printf("This is unused.\n")

	// Let the compiler know that `str` is "used" by the program.
	_ = str
}
