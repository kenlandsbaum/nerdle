package pro

func GetPerson() *Person {
	person := Person{
		Id:        1,
		FirstName: "ken",
		LastName:  "L",
		Email:     "ken@mail.com",
		Phone:     []*Person_PhoneNumber{{Number: "1234567890", Type: Person_MOBILE}},
		Address:   &Address{Street: "123 elm st", State: "TX", Zipcode: "75000"},
	}
	return &person
}
