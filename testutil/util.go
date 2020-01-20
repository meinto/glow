package testutil

import "reflect"

func ForEachTestSet(testSet interface{}, cb func(interface{})) {
	switch reflect.TypeOf(testSet).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(testSet)

		for i := 0; i < s.Len(); i++ {
			cb(s.Index(i).Interface())
		}
	}
}
