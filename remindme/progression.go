package remindme

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
)

func newProgressionWithWorkingDirectory() *progression {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Unrecoverable error fetching working directory")
	}
	return newProgression(path.Join(wd, "remindme_progress.json"))
}

func newProgression(stateFilePath string) *progression {
	return &progression{stateFilePath: stateFilePath}
}

type progression struct {
	stateFilePath string
}

var errExhaustedIndex = errors.New("ExhaustedIndexError")

type progFileState struct {
	TopicToIndexMap map[string]int
}

func (p *progression) nextIndexForTopic(topic string, limit int) (int, error) {
	_, err := os.Stat(p.stateFilePath)

	if errors.Is(err, os.ErrNotExist) {
		createDefaultStateFile(p.stateFilePath)
	}

	state := loadStateFile(p.stateFilePath)

	index := state.TopicToIndexMap["example"]
	if index >= limit {
		return -1, errExhaustedIndex
	}
	state.TopicToIndexMap["example"] = index + 1
	saveStateFile(p.stateFilePath, state)

	return index, nil
}

func createDefaultStateFile(path string) {
	saveStateFile(path, &progFileState{TopicToIndexMap: map[string]int{}})
}

func saveStateFile(path string, state *progFileState) {
	data, err := json.Marshal(state)
	if err != nil {
		log.Fatal("Unrecoverable JSON error: ", err)
	}

	err = os.WriteFile(path, data, 0644)

	if err != nil {
		log.Fatal("Unrecoverable error while saving state: ", err)
	}
}

func loadStateFile(path string) *progFileState {
	data, err := os.ReadFile(path)

	if err != nil {
		log.Fatal("Unrecoverable progression file load error: ", err)
	}
	state := &progFileState{}
	json.Unmarshal(data, state)

	return state
}
