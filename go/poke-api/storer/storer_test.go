package storer

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"testing"

	"playground/poke-api/types"

	"github.com/jackc/pgx/v5/pgxpool"
	"gotest.tools/v3/assert"
)

func setup(ctx context.Context) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		log.Panic("connecting to db", err)
	}
	// Drop table if exists (for repeatable tests)
	_, _ = pool.Exec(ctx, `DROP TABLE IF EXISTS pokedex;`)

	// Create table
	_, err = pool.Exec(ctx, `
		CREATE TABLE pokedex (
			id SERIAL PRIMARY KEY,
			entry INT,
			pokemon JSONB NOT NULL
		);
	`)
	if err != nil {
		log.Panic("creating table", err)
	}
	return pool
}

func TestGetPokemon(t *testing.T) {
	ctx := context.Background()
	p := setup(ctx)
	defer p.Close()

	expectedPokemon := &types.Pokemon{
		Id:        1,
		Name:      "Bulbasaur",
		Type1:     "Grass",
		Type2:     "Poison",
		Abilities: [2]string{"Overgrow", "Chlorophyll"},
		ImageUrl:  "bulbasaur.png",
	}

	d := Database{
		pool: p,
		log:  slog.Default(),
	}
	pokemonJSON, err := json.Marshal(expectedPokemon)
	_, err = d.pool.Exec(ctx,
		`INSERT INTO pokedex (entry, pokemon) VALUES ($1, $2)`, 1, pokemonJSON,
	)
	assert.NilError(t, err, "failed to insert test pokemon")

	got, err := d.GetPokemon(ctx, 1)
	assert.NilError(t, err)
	assert.DeepEqual(t, got, expectedPokemon)
}
