package types // nolint: dupl

import (
	"github.com/FishtechCSOC/locomotive/v8/pkg/helpers"
)

type Customer struct {
	Name string `json:"name" mapstructure:"name"`
	ID   string `json:"id" mapstructure:"id"`
}

func (customer *Customer) DeepCopy() Customer {
	return Customer{
		Name: customer.Name,
		ID:   customer.ID,
	}
}

func (customer *Customer) MergeLeft(other Customer) Customer {
	return Customer{
		Name: helpers.DefaultString(other.Name, customer.Name),
		ID:   helpers.DefaultString(other.ID, customer.ID),
	}
}

func (customer *Customer) MergeRight(other Customer) Customer {
	return Customer{
		Name: helpers.DefaultString(customer.Name, other.Name),
		ID:   helpers.DefaultString(customer.ID, other.ID),
	}
}

func (customer *Customer) Empty() bool {
	return customer.Name == "" && customer.ID == ""
}
