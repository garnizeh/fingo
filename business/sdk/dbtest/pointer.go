package dbtest

import (
	"github.com/garnizeh/fingo/business/types/money"
	"github.com/garnizeh/fingo/business/types/name"
	"github.com/garnizeh/fingo/business/types/password"
	"github.com/garnizeh/fingo/business/types/quantity"
)

// NamePointer is a helper to get a *Name from a string. It's in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func NamePointer(value string) *name.Name {
	name := name.MustParse(value)
	return &name
}

// NameNullPointer is a helper to get a *EmptyName from a string. It's in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func NameNullPointer(value string) *name.Null {
	name := name.MustParseNull(value)
	return &name
}

// MoneyPointer is a helper to get a *Money from a float. It's in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func MoneyPointer(value float64) *money.Money {
	money := money.MustParse(value)
	return &money
}

// QuantityPointer is a helper to get a *Quantity from an int. It's in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func QuantityPointer(value int) *quantity.Quantity {
	quantity := quantity.MustParse(value)
	return &quantity
}

// PasswordPointer is a helper to get a *Password from a string. It's in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func PasswordPointer(value string) *password.Password {
	pass := password.MustParse(value)
	return &pass
}
