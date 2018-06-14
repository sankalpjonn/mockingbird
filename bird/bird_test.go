package bird

import (
	"reflect"
	"testing"
)

var (
	testEgg = &Egg{
		Id: "test",
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
	b.createEgg(testEgg)
}

func TestGetEgg(t *testing.T) {
	b := New()
	err, egg := b.getEgg(testEgg.Id)
	if err != nil {
		panic(err)
	}
	egg.TTL = testEgg.TTL
	if !reflect.DeepEqual(egg, testEgg) {
		t.Error("Eggs do not match")
	}
}
