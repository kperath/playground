package pokedex

import (
	"context"

	"playground/poke-api/storer"
	"playground/poke-api/types"
)

type PaginatedPokemon struct {
	Results    []*types.Pokemon `json:"results"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalCount int64            `json:"total_count"`
}

type Pokedex struct {
	*public
	*admin
}

type public struct {
	db     *storer.Database
	search *storer.Search
}

type admin struct {
	db     *storer.Database
	search *storer.Search
}

func New(db *storer.Database, search *storer.Search) *Pokedex {
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

func (p *public) GetPokemon(ctx context.Context, pokedexEntry int) (*types.Pokemon, error) {
	return p.db.GetPokemon(ctx, pokedexEntry)
}

type SearchOpts struct {
	Page     int
	PageSize int
}

func defaultSearchOpts() *SearchOpts {
	return &SearchOpts{
		Page:     1,
		PageSize: 10,
	}
}

func (p *public) SearchPokemon(ctx context.Context, name string, opts *SearchOpts) (*PaginatedPokemon, error) {
	if opts == nil {
		opts = defaultSearchOpts()
	}

	pkmns, resultCount, err := p.search.SearchPokemon(ctx, name, opts.Page, opts.PageSize)
	if err != nil {
		return nil, err
	}
	return &PaginatedPokemon{
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalCount: resultCount,
		Results:    pkmns,
	}, nil
}

func (p *public) AddPokemon(ctx context.Context, pkmn *types.Pokemon) error {
	return p.db.AddPokemon(ctx, pkmn)
}

func (p *admin) DeletePokemon(ctx context.Context, pokedexEntry int) error {
	return p.db.DeletePokemon(ctx, pokedexEntry)
}

func (p *admin) UpdatePokemon(ctx context.Context, pokedexEntry int, pkmn *types.Pokemon) (*types.Pokemon, error) {
	return p.db.UpdatePokemon(ctx, pokedexEntry, pkmn)
}
