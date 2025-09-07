package pokemon

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	testPokemon = &Pokemon{
		PokedexNumber:    1,
		IsLegendary:      false,
		MegaEvolved:      false,
		Generation:       "1",
		NickName:         "Bulba",
		Name:             "Bulbasaur",
		JapaneseName:     "Fushigidaneフシギダネ",
		Classification:   "Seed Pokémon",
		Type1:            "grass",
		Type2:            "poison",
		Abilities:        [2]string{"Overgrow", "Chlorophyll"},
		ExperienceGrowth: 1059860,
		Experience:       52993,
		BaseTotal:        318,
		Attack:           49,
		Defense:          49,
		Speed:            45,
		SpAttack:         65,
		SpDefense:        65,
		HP:               45,
		Level:            5,
		Item:             "Oran Berry",
		Nature:           "Adament",
		Attacks:          [4]string{"Tackle", "", "", ""},
		Image:            []byte{0x42, 0x75, 0x6c, 0x62, 0x61},
	}

	testPokemonData = `[{
	"pokedex_number": 1,
	"is_legendary": false,
	"mega_evolved": false,
	"generation": "1",
	"nick_name": "Bulba",
	"name": "Bulbasaur",
	"japanese_name": "Fushigidaneフシギダネ",
	"classification": "Seed Pokémon",
	"type1": "grass",
	"type2": "poison",
	"abilities": ["Overgrow", "Chlorophyll"],
	"experience_growth": 1059860,
	"experience": 52993,
	"base_total": 318,
	"attack": 49,
	"defense": 49,
	"speed": 45,
	"sp_attack": 65,
	"sp_defense": 65,
	"hp": 45,
	"level": 5,
	"item": "Oran Berry",
	"nature": "Adament",
	"attacks": ["Tackle", "", "", ""],
	"image": "QnVsYmE="
}]`
)

func TestNewParty(t *testing.T) {
	f := strings.NewReader(testPokemonData)
	dex, err := New(f)
	if err != nil {
		t.Fatal(err)
	}

	got := dex.Pokemon[0]
	expect := testPokemon

	equal := cmp.Diff(got, expect, cmpopts.SortSlices(func(a, b string) bool {
		return a < b
	}))
	if equal != "" {
		t.Fatal(equal)
	}
}
