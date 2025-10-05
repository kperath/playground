package pokedex

import (
	"playground/poke-api/storer"
	"playground/poke-api/types"
)

type PaginatedPokemon struct {
	Results    []*types.Pokemon `json:"results"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalCount int              `json:"total_count"`
}

type Pokedex struct {
	*public
	*admin
}

type public struct {
	db     *storer.Database
	search *storer.Searcher
}

type admin struct {
	db     *storer.Database
	search *storer.Searcher
}

func New(db *storer.Database, search *storer.Searcher) *Pokedex {
	return &Pokedex{
		public: &public{
			db:     db,
			search: search,
		},
		admin: &admin{
			db:     db,
			search: search,
		},
	}
}

func (p *public) GetPokemon() ([]*PaginatedPokemon, error) {
	return nil, nil
}

func (p *public) SearchPokemon() ([]*PaginatedPokemon, error) {
	return nil, nil
}

func (p *public) AddPokemon() error {
	return nil
}

func (p *admin) DeletePokemon() error {
	return nil
}

func (p *admin) UpdatePokemon() (*PaginatedPokemon, error) {
	return nil, nil
}
