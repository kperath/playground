package storer

import (
	"context"
	"encoding/json"
	"log/slog"

	"playground/poke-api/types"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

type Searcher struct {
	log *slog.Logger
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		pool: pool,
		log:  slog.Default(),
	}
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

func (s *Searcher) SearchPokemon() ([]*types.Pokemon, error) {
	return nil, nil
}
