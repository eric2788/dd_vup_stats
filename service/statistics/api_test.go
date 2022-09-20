package statistics

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	res, err := GetUserInfo(2299184)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, int64(2299184), res.Data.Mid)
	assert.NotEqual(t, 0, res.Data.Official.Role)
	res, err = GetUserInfo(85650700)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, int64(85650700), res.Data.Mid)
	assert.Equal(t, 0, res.Data.Official.Role)
}
