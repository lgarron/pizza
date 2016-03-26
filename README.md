# Unexpected behaviour of `go vet -unusedresult`

## Issue Summary

There are 4 functions in `pizza.go` and 4 functions in `string.go` that `go vet` cannot report as an unused result as of [21cc49bd](https://github.com/golang/tools/blob/21cc49bd030cf5c6ebaca2fa0e3323628efed6d8/cmd/vet/unused.go) (March 25, 2016).

As an example, the functions in `pizza.go` all share the same implementation:

    return Pizza{
      toppings: append(pizza.toppings, topping),
    }

... although they have four different signatures:

    func (pizza Pizza) addToppingMethod(topping string) Pizza
    func (pizza Pizza) PublicAddToppingMethod(topping string) Pizza
    func addToppingBareFunction(pizza Pizza, topping string) Pizza
    func PublicAddToppingBareFunction(pizza Pizza, topping string) Pizza

It is not possible to check for unused results of these functions using `go vet` with the `-unusedresult` flag.

Based on [the commit that introduced `-unusedresult`](https://github.com/golang/tools/commit/4a08fb6fc335489d063834778b71d1af382f4c27), it seems this was originally designed to catch "pure" functions. However, [the documentation](https://golang.org/cmd/vet/#hdr-Unused_result_of_certain_function_calls) does not make this clear at all.

Also:

- It's not obvious which function types `-unusedresult` can check for, which functions to pass to `-unusedfuncs` or `-unusedstringmethods`, and how to format and qualify them.
- It's not possible to get `go vet` to output the default values for `-unusedfuncs` and `-unusedstringmethods` and the documentation is vague ("By default, this includes functions like fmt.Errorf and fmt.Sprintf and methods like String and Error."). In addition, there is no way to use the default values *in addition to* custom values in a single invocation, which makes it easy to miss out on vetting the default values *by accident*. The only way to do this seems to be to run `go vet` multiple times` with different arguments.

## Reproduction Steps

    # Get the code
    go get github.com/lgarron/pizza
    cd "$GOPATH/src/github.com/lgarron/pizza"

    # Get `go vet`
    go get golang.org/x/tools/cmd/vet

    # Sanity check that the code compiles and runs.
    go test -v

    # Run `go vet` with default settings.
    go tool vet -unusedresult *_test.go

    # Run `go vet` using qualified *and* unqualified versions for each function.
    go tool vet -unusedresult -unusedfuncs="makeStringMethod,PublicMakeStringMethod,makeStringBareFunction,PublicMakeStringBareFunction,addToppingMethod,PublicAddToppingMethod,addToppingBareFunction,PublicAddToppingBareFunction,pizza.Pizza.makeStringMethod,pizza.Pizza.PublicMakeStringMethod,pizza.makeStringBareFunction,pizza.PublicMakeStringBareFunction,pizza.Pizza.addToppingMethod,pizza.Pizza.PublicAddToppingMethod,pizza.addToppingBareFunction,pizza.PublicAddToppingBareFunction,errors.New,fmt.Errorf,fmt.Sprintf,fmt.Sprint,sort.Reverse" -unusedstringmethods="makeStringMethod,PublicMakeStringMethod,makeStringBareFunction,PublicMakeStringBareFunction,addToppingMethod,PublicAddToppingMethod,addToppingBareFunction,PublicAddToppingBareFunction,pizza.Pizza.makeStringMethod,pizza.Pizza.PublicMakeStringMethod,pizza.makeStringBareFunction,pizza.PublicMakeStringBareFunction,pizza.Pizza.addToppingMethod,pizza.Pizza.PublicAddToppingMethod,pizza.addToppingBareFunction,pizza.PublicAddToppingBareFunction,Error,String" *_test.go

### Debugging using manual traces:

The relevant file int he `go vet` source is called `unused.go`. We will use [a modified version of `unused.go`](./modified-vet-src/unused.go) to show us what is happening.

    # Check out a known-bad version of the source for `go vet`.
    cd "$GOPATH/src/golang.org/x/tools"
    git checkout 21cc49bd030cf5c6ebaca2fa0e3323628efed6d8
    cd "$GOPATH/src/github.com/lgarron/pizza"

    # Overwrite the relevant `go vet` code with an annotated version.
    cp modified-vet-src/unused.go "$GOPATH/src/golang.org/x/tools/cmd/vet/unused.go"

    # Build a new version of `vet` in the current project directory.
    go build -a golang.org/x/tools/cmd/vet

    # View output.
    ./vet -unusedresult *_test.go

Here is the output I get:

    unusedFuncsFlag: errors.New,fmt.Errorf,fmt.Sprintf,fmt.Sprint,sort.Reverse
    unusedStringMethodsFlag: Error,String

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [checking for unusedFuncs]
    REPORT: qualified name `fmt.Sprintf` is in unusedFuncs

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [pkg.uses was not okay]

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [pkg.uses was not okay]

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    NO REPORT: not method + unqualified

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    NO REPORT: not method + unqualified

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [checking for unusedFuncs]
    REPORT: qualified name `fmt.Sprintf` is in unusedFuncs

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [checking for unusedFuncs]
    REPORT: qualified name `fmt.Sprintf` is in unusedFuncs

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [pkg.uses was not okay]

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [pkg.uses was not okay]

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    NO REPORT: not method + unqualified

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    NO REPORT: not method + unqualified

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [checking for unusedFuncs]
    REPORT: qualified name `fmt.Sprintf` is in unusedFuncs

    [start]
    [checking for conversion]
    [checking for (not method + unqualified)]
    [checking pkg.selectors]
    [pkg.selectors not okay]
    [checking pkg.uses]
    [checking for unusedFuncs]
    NO REPORT BUT REPORTABLE: qualified name `fmt.Printf` is not in unusedFuncs

## Suggestions

- 1a) Remove some of the restrictions and assumptions in the source to allow reporting these functions as unused, or 1b) explicitly document the limited behaviour.
- 2) Document which functions can be passed in, with clear examples of how to format them.

I tried 1a), but I don't know the code well enough to contribute a fix yet.
