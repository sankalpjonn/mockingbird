package bird

import (
	"reflect"
	"testing"
)

var (
	testEggId = "test"
	testEgg   = &Egg{
		Headers: map[string]string{
			"Content-Type":    "application/json",
			"X-Forwarded-For": "1.1.1.1",
		},
		Body:       "awesome!",
		StatusCode: 204,
		TTL:        10,
	}
)

func TestCreateEgg(*testing.T) {
	b := New()
	err := b.createEgg(testEggId, testEgg)
	if err != nil {
		panic(err)
	}
}

func TestGetEgg(t *testing.T) {
	b := New()
	err, egg := b.getEgg(testEggId)
	if err != nil {
		panic(err)
	}
	egg.TTL = testEgg.TTL
	if !reflect.DeepEqual(egg, testEgg) {
		t.Error("Eggs do not match")
	}
}
