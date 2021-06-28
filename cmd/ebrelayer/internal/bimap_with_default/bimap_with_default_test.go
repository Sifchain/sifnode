package bimap_with_default

import (
	"github.com/stretchr/testify/assert"
	"github.com/vishalkuo/bimap"
	"testing"
)

func TestNewBiMapWithDefault(t *testing.T) {
	dflt := "some default"
	objectUnderTest := bimap.NewBiMap()
	objectUnderTest.Insert("one", 1)
	assert.Equal(t, GetWithDefault(objectUnderTest, "one", dflt), 1)
	assert.Equal(t, GetInverseWithDefault(objectUnderTest, 1, dflt), "one")
	assert.Equal(t, GetWithDefault(objectUnderTest, "missing", dflt), dflt)
	assert.Equal(t, GetInverseWithDefault(objectUnderTest, 2, dflt), dflt)
}
