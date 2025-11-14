package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidNewStreamId(t *testing.T) {
	s := NewStream()
	isValid, err := s.IsNewStreamIdValid("1", "0")
	assert.True(t, isValid)
	assert.Nil(t, err)

	// add these keys to the stream
	s.AddIdWithKV("1", "0", map[string]string{"foo": "bar"})

	// rerun tests
	isValid, err = s.IsNewStreamIdValid("1", "0")
	assert.Nil(t, err)
	assert.False(t, isValid)
}

func TestXrangeReturn_RArray(t *testing.T) {
	id := "1526985054069-0"
	k1 := "temperature"
	v1 := "36"
	k2 := "humidity"
	v2 := "95"
	r := XRangeReturn{
		ID: id,
		KV: map[string]string{
			k1: v1,
			k2: v2,
		},
	}

	rArray := r.ToRArray()
	assert.Equal(t, "1526985054069-0", rArray[0].(string))

	innerArray := rArray[1].([]any)
	valMap := make(map[string]string)
	prevKey := ""
	for idx, val := range innerArray {
		if idx%2 == 0 {
			prevKey = val.(string)
			continue
		}
		valMap[prevKey] = val.(string)
	}
	assert.Equal(t, v1, valMap[k1])
	assert.Equal(t, v2, valMap[k2])
}
