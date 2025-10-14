package pokedex

import (
	"context"
	"log/slog"

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
	log    *slog.Logger
}

type admin struct {
	db     *storer.Database
	search *storer.Search
	log    *slog.Logger
}

func New(db *storer.Database, search *storer.Search) *Pokedex {
	return &Pokedex{
		public: &public{
			db:     db,
			search: search,
			log:    slog.Default(),
		},
		admin: &admin{
			db:     db,
			search: search,
			log:    slog.Default(),
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
	if pkmn.Id == 0 {
		p.log.Error("adding: empty pokedex entry value")
		return types.ErrAddPokemon
	}
	err := p.db.AddPokemon(ctx, pkmn)
	if err != nil {
		p.log.Error("adding pokemon", "db", "postgres")
		return err
	}
	err = p.search.AddPokemon(ctx, pkmn)
	if err != nil {
		p.log.Error("adding pokemon", "db", "elasticsearch")
		return err
	}
	return nil
}

func (a *admin) DeletePokemon(ctx context.Context, pokedexEntry int) error {
	return a.db.DeletePokemon(ctx, pokedexEntry)
}

func (a *admin) UpdatePokemon(ctx context.Context, pokedexEntry int, pkmn *types.Pokemon) (*types.Pokemon, error) {
	if pokedexEntry == 0 {
		a.log.Error("updating: empty pokedex entry value")
		return nil, types.ErrUpdatePokemon
	}
	return a.db.UpdatePokemon(ctx, pokedexEntry, pkmn)
}
