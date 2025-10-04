package pokemon

import (
	"encoding/json"
	"io"
)

// Pokemon ...
type Pokemon struct {
	PokedexNumber    int       `json:"pokedex_number,omitempty"`
	IsLegendary      bool      `json:"is_legendary,omitempty"`
	MegaEvolved      bool      `json:"mega_evolved,omitempty"`
	Generation       string    `json:"generation,omitempty"`
	NickName         string    `json:"nick_name,omitempty"`
	Name             string    `json:"name,omitempty"`
	JapaneseName     string    `json:"japanese_name,omitempty"`
	Classification   string    `json:"classification,omitempty"`
	Type1            string    `json:"type1,omitempty"`
	Type2            string    `json:"type2,omitempty"`
	Abilities        [2]string `json:"abilities,omitempty"`
	ExperienceGrowth int       `json:"experience_growth,omitempty"`
	Experience       int       `json:"experience,omitempty"`
	BaseTotal        int       `json:"base_total,omitempty"`
	Attack           int       `json:"attack,omitempty"`
	Defense          int       `json:"defense,omitempty"`
	Speed            int       `json:"speed,omitempty"`
	SpAttack         int       `json:"sp_attack,omitempty"`
	SpDefense        int       `json:"sp_defense,omitempty"`
	HP               int       `json:"hp,omitempty"`
	Level            int       `json:"level,omitempty"`
	Item             string    `json:"item,omitempty"`
	Nature           string    `json:"nature,omitempty"`
	Attacks          [4]string `json:"attacks,omitempty"`
	Image            []byte    `json:"image,omitempty"`
}

// Party ...
type Party struct {
	Pokemon []*Pokemon `json:"pokemon,omitempty"`
}

// NewParty inits the client
func NewParty(data io.Reader) (*Party, error) {
	var pokemon []*Pokemon
	err := json.NewDecoder(data).Decode(&pokemon)
	if err != nil {
		return nil, err
	}
	return &Party{
		Pokemon: pokemon,
	}, nil
}
