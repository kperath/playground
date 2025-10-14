package types

import "errors"

// Pokemon ...
type Pokemon struct {
	Id        int      `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Type1     string   `json:"type1,omitempty"`
	Type2     string   `json:"type2,omitempty"`
	Abilities []string `json:"abilities,omitempty"`
	ImageUrl  string   `json:"image,omitempty"`
}

var ErrGetPokemon = errors.New("failed to get pokemon")
var ErrAddPokemon = errors.New("failed to add pokemon")
var ErrDeletePokemon = errors.New("failed to delete pokemon")
var ErrUpdatePokemon = errors.New("failed to update pokemon")
var ErrSearchPokemon = errors.New("failed to search pokemon")
