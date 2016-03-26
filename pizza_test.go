package pizza

import (
	"fmt"
	"testing"
)

func TestPizza(t *testing.T) {
	var p Pizza
	var str string

	// We place known-reported calls at the start and end so we can tell
	// `go vet` is actually looking for them.
	str = fmt.Sprintf("This is used.")
	fmt.Sprintf("This is unused.")

	p = p.addToppingMethod("used private method topping")
	p.addToppingMethod("unused private method topping")

	p = p.PublicAddToppingMethod("used public method topping")
	p.PublicAddToppingMethod("unused public method topping")

	p = addToppingBareFunction(p, "used locally-unqualified private bare call topping")
	addToppingBareFunction(p, "unused locally-unqualified private bare call topping")

	p = PublicAddToppingBareFunction(p, "used locally-unqualified public bare call topping")
	PublicAddToppingBareFunction(p, "unused locally-unqualified public bare call topping")

	// Can't qualify an identifier in this package itself.
	//
	// p = pizza.combineWithTopping(p, "used package-qualified private bare call topping")
	// pizza.combineWithTopping(p, "unused package-qualified private bare call topping")
	//
	// p = pizza.PublicCombineWithTopping(p, "used package-qualified public bare call topping")
	// pizza.PublicCombineWithTopping(p, "unused package-qualified public bare call topping")

	str = fmt.Sprintf("This is printed and the result is used.")
	fmt.Sprintf("This is printed and the result is unused.")

	// Let the compiler know that `str` is "used" by the program.
	_ = str
}
