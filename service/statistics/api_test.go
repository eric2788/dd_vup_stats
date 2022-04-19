package statistics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	res, err := GetUserInfo(2299184)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, res.Code, 0)
	assert.Equal(t, res.Data.Mid, int64(2299184))
	assert.NotEqual(t, res.Data.Official.Role, 0)

	res, err = GetUserInfo(85650700)
	assert.Equal(t, res.Code, 0)
	assert.Equal(t, res.Data.Mid, int64(85650700))
	assert.Equal(t, res.Data.Official.Role, 0)
}
