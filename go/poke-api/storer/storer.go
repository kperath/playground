package storer

import (
	"context"
	"log/slog"

	"playground/poke-api/types"

	"github.com/jackc/pgx/v5/pgtype"
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

//TODO: api, and instrumentation

func (db *Database) GetPokemon(ctx context.Context, pokedexNumber int) (*types.Pokemon, error) {
	var pokemonJSON pgtype.JSONBCodec
	err := db.pool.QueryRow(ctx, `
		SELECT pokemon FROM pokedex WHERE pokedexNumber = $1
		`, pokedexNumber).Scan(&pokemonJSON)
	if err != nil {
		db.log.Error("getting pokemon", "error", err.Error())
		return nil, types.ErrGetPokemon
	}
	pkmn := &types.Pokemon{}
	pokemonJSON.Unmarshal()
	// if err := json.Unmarshal()
	return pkmn, nil
}
func (db *Database) AddPokemon() error {
	return nil
}
func (db *Database) DeletePokemon() error {
	return nil
}
func (db *Database) UpdatePokemon() (*types.Pokemon, error) {
	return nil, nil
}

func (s *Searcher) SearchPokemon() ([]*types.Pokemon, error) {
	return nil, nil
}
