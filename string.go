package pizza

func (pizza Pizza) makeStringMethod() string {
	return "string from private method"
}

func (pizza Pizza) PublicMakeStringMethod() string {
	return "string from public method"
}

func makeStringBareFunction() string {
	return "string from private bare function"
}

func PublicMakeStringBareFunction() string {
	return "string from private bare function"
}
