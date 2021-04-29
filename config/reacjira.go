package config

import (
	"log"

	toml "github.com/sioncojp/tomlssm"
	"golang.org/x/xerrors"
)

type Reacjiras struct {
	Values []Reacjira `toml:"reacjiras"`
}

type Reacjira struct {
	Emoji       string `toml:"emoji"`
	Project     string `toml:"project"`
	IssueType   string `toml:"issue_type"`
	EpicKey     string `toml:"epic_key"`
	Description string `toml:"description"`
}

func LoadReacjiraToml(filename string) (*Reacjiras, error) {
	log.Printf("try to load reacjira from %s", filename)

	var reacjiras Reacjiras
	if _, err := toml.DecodeFile(filename, &reacjiras, "ap-northeast-1"); err != nil {
		return nil, xerrors.Errorf("failed to load reacjira toml from %s: %w", filename, err)
	}

	return &reacjiras, nil
}
