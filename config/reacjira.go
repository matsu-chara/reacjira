package config

import (
	"fmt"
	"log"

	toml "github.com/sioncojp/tomlssm"
)

type reacjiras struct {
	Reacjiras []Reacjira `toml:"reacjiras"`
}

type Reacjira struct {
	Emoji        string `toml:"emoji"`
	Project      string `toml:"project"`
	IssueType    string `toml:"issue_type"`
	EpicKey      string `toml:"epic_key"`
	Description  string `toml:"description"`
	SlackMessage string `toml:"slack_message"`
}

func LoadReacjiraToml(filename string) ([]Reacjira, error) {
	log.Printf("try to load reacjira from %s", filename)

	var reacjiras reacjiras
	if _, err := toml.DecodeFile(filename, &reacjiras, "ap-northeast-1"); err != nil {
		return nil, fmt.Errorf("failed to load reacjira toml from %s: %w", filename, err)
	}

	return reacjiras.Reacjiras, nil
}
