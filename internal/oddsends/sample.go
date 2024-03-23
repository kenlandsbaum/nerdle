package oddsends

import "reflect"

func CheckDoesImplement[T any](arg any) bool {
	genType := reflect.TypeOf(new(T)).Elem()
	tp := reflect.TypeOf(arg)

	return tp.Implements(genType)
}
