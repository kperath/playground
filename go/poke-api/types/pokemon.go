package types

import "errors"

// Pokemon ...
type Pokemon struct {
	Id        int       `json:"pokedex_number,omitempty"`
	Name      string    `json:"name,omitempty"`
	Type1     string    `json:"type1,omitempty"`
	Type2     string    `json:"type2,omitempty"`
	Abilities [2]string `json:"abilities,omitempty"`
	ImageUrl  string    `json:"image,omitempty"`
}

var ErrGetPokemon = errors.New("failed to get pokemon")
