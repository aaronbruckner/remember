package remindme

import (
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newStateFilePath() string {
	return path.Join(os.TempDir(), uuid.NewString()+".json")
}

func newDefaultProgression() *progression {
	return newProgression(newStateFilePath())
}

func Test_nextIndexForTopic_returnsIncrementalIndex(t *testing.T) {
	// Given
	progression := newDefaultProgression()
	assert := assert.New(t)
	var indexes []int

	// When
	for i := 0; i < 5; i++ {
		index, err := progression.nextIndexForTopic("example", 5)
		assert.Nil(err)
		indexes = append(indexes, index)
	}

	// Then
	assert.Equal([]int{0, 1, 2, 3, 4}, indexes)
}

func Test_nextIndexForTopic_returnsNilIfLimitReached(t *testing.T) {
	// Given
	progression := newDefaultProgression()
	assert := assert.New(t)

	// When
	progression.nextIndexForTopic("example", 2)
	progression.nextIndexForTopic("example", 2)

	index1, err1 := progression.nextIndexForTopic("example", 2)
	index2, err2 := progression.nextIndexForTopic("example", 2)

	// Then
	assert.ErrorIs(err1, ErrExhaustedIndex)
	assert.ErrorIs(err2, ErrExhaustedIndex)
	assert.Equal(-1, index1)
	assert.Equal(-1, index2)
}

func Test_nextIndexForTopic_keepsTrackOfReturnedIndexBetweenInstances(t *testing.T) {
	// Given
	progression1 := newDefaultProgression()
	progression2 := newProgression(progression1.stateFilePath)
	assert := assert.New(t)

	// When
	index1, err1 := progression1.nextIndexForTopic("example", 2)
	index2, err2 := progression2.nextIndexForTopic("example", 2)
	index3, err3 := progression1.nextIndexForTopic("example", 2)

	// Then
	assert.Nil(err1)
	assert.Nil(err2)
	assert.ErrorIs(err3, ErrExhaustedIndex)

	assert.Equal(0, index1)
	assert.Equal(1, index2)
	assert.Equal(-1, index3)
}
