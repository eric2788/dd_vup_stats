package statistics

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetUserInfo(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	res, err := GetUserInfoRetry(2299184, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, int64(2299184), res.Data.Mid)
	assert.NotEqual(t, 0, res.Data.Official.Role)
	res, err = GetUserInfoRetry(85650700, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, int64(85650700), res.Data.Mid)
	assert.Equal(t, 0, res.Data.Official.Role)
}

func TestGetVtbMoe(t *testing.T) {
	vtbs, err := GetVtbListVtbMoe()
	if err != nil {
		t.Skip(err)
	}
	t.Log(vtbs)
}

func TestGetVtbOoo(t *testing.T) {
	vtbs, err := GetVtbListLaplace()
	if err != nil {
		t.Skip(err)
	}
	t.Log(vtbs)
}
