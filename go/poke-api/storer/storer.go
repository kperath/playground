package storer

import (
	"context"
	"encoding/json"
	"log/slog"

	"playground/poke-api/types"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	estypes "github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

type Search struct {
	log    *slog.Logger
	client *elasticsearch.TypedClient
}

func (db *Database) GetPokemon(ctx context.Context, pokedexEntry int) (*types.Pokemon, error) {
	var pokemonJSON []byte
	err := db.pool.QueryRow(ctx, `
		SELECT pokemon FROM pokedex WHERE entry = $1
		`, pokedexEntry).Scan(&pokemonJSON)
	if err != nil {
		db.log.Error("getting pokemon", "error", err.Error())
		return nil, types.ErrGetPokemon
	}
	pkmn := &types.Pokemon{}
	json.Unmarshal(pokemonJSON, pkmn)
	return pkmn, nil
}

func (db *Database) AddPokemon(ctx context.Context, pkmn *types.Pokemon) error {
	_, err := db.pool.Exec(ctx, `
		INSERT INTO pokedex(entry,pokemon)  VALUES ($1, $2)
	`, pkmn.Id, pkmn)
	if err != nil {
		db.log.Error("adding pokemon", "error", err.Error())
		return types.ErrAddPokemon
	}
	return nil
}

func (db *Database) DeletePokemon(ctx context.Context, pokedexEntry int) error {
	_, err := db.pool.Exec(ctx, `
		DELETE FROM pokedex WHERE entry = $1
		`, pokedexEntry)
	if err != nil {
		db.log.Error("deleting pokemon", "error", err.Error())
		return types.ErrDeletePokemon
	}
	return nil
}

func (db *Database) UpdatePokemon(ctx context.Context, pokedexEntry int, pkmn *types.Pokemon) (*types.Pokemon, error) {
	updatedPokemon := &types.Pokemon{}
	err := db.pool.QueryRow(ctx, `
		UPDATE pokedex
		SET pokemon = $1
		WHERE entry = $2
		RETURNING pokemon
		`, pkmn, pokedexEntry).Scan(&updatedPokemon)
	if err != nil {
		db.log.Error("deleting pokemon", "error", err.Error())
		return nil, types.ErrUpdatePokemon
	}
	return updatedPokemon, nil
}

func (s *Search) SearchPokemon(ctx context.Context, name string, page, pageSize int) ([]*types.Pokemon, int64, error) {
	from := (page - 1) * pageSize
	size := pageSize
	res, err := s.client.
		Search().
		Index("pokedex").
		From(from).
		Size(size).
		Request(&search.Request{
			Query: &estypes.Query{
				Match: map[string]estypes.MatchQuery{
					"name": {
						Query: name,
					},
				},
			},
		}).Do(ctx)
	if err != nil {
		s.log.Error("searching pokemon", "error", err)
		return nil, 0, types.ErrSearchPokemon
	}
	var pkmns []*types.Pokemon
	for _, pokemon := range res.Hits.Hits {
		p := &types.Pokemon{}
		err := json.Unmarshal(pokemon.Source_, &p)
		if err != nil {
			s.log.Error("decoding pokemon", "error", err)
			return nil, 0, types.ErrSearchPokemon
		}
		pkmns = append(pkmns, p)
	}
	resultCount := res.Hits.Total.Value
	return pkmns, resultCount, nil
}
