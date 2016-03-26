package pizza

type Pizza struct {
	toppings []string
}

func (pizza Pizza) addToppingMethod(topping string) Pizza {
	return Pizza{
		toppings: append(pizza.toppings, topping),
	}
}

func (pizza Pizza) PublicAddToppingMethod(topping string) Pizza {
	return Pizza{
		toppings: append(pizza.toppings, topping),
	}
}

func addToppingBareFunction(pizza Pizza, topping string) Pizza {
	return Pizza{
		toppings: append(pizza.toppings, topping),
	}
}

func PublicAddToppingBareFunction(pizza Pizza, topping string) Pizza {
	return Pizza{
		toppings: append(pizza.toppings, topping),
	}
}
