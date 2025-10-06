package storer

import (
	"context"
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

func cleanup(ctx context.Context, pool *pgxpool.Pool) {
	defer pool.Close()
	_, err := pool.Exec(ctx, `DROP TABLE IF EXISTS pokedex;`)
	if err != nil {
		log.Panic(err)
	}
}

func testPokemon() *types.Pokemon {
	return &types.Pokemon{
		Id:        1,
		Name:      "Bulbasaur",
		Type1:     "Grass",
		Type2:     "Poison",
		Abilities: [2]string{"Overgrow", "Chlorophyll"},
		ImageUrl:  "bulbasaur.png",
	}
}

func TestGetPokemon(t *testing.T) {
	ctx := context.Background()
	p := setup(ctx)
	defer cleanup(ctx, p)

	expectedPokemon := testPokemon()
	d := Database{
		pool: p,
		log:  slog.Default(),
	}
	_, err := d.pool.Exec(ctx,
		`INSERT INTO pokedex (entry, pokemon) VALUES ($1, $2)`, 1, expectedPokemon,
	)
	assert.NilError(t, err, "failed to insert test pokemon")

	got, err := d.GetPokemon(ctx, 1)
	assert.NilError(t, err)
	assert.DeepEqual(t, got, expectedPokemon)
}

func TestAddPokemon(t *testing.T) {
	ctx := context.Background()
	p := setup(ctx)
	defer cleanup(ctx, p)

	expectedPokemon := testPokemon()
	d := Database{
		pool: p,
		log:  slog.Default(),
	}
	err := d.AddPokemon(ctx, expectedPokemon)
	assert.NilError(t, err)

	gotPokemon := &types.Pokemon{}
	err = d.pool.QueryRow(ctx, `
		SELECT pokemon
		FROM pokedex
		WHERE id = $1
	`, 1).Scan(&gotPokemon)
	assert.NilError(t, err)
	assert.DeepEqual(t, gotPokemon, expectedPokemon)
}

func TestDeletePokemon(t *testing.T) {
	ctx := context.Background()
	p := setup(ctx)
	defer cleanup(ctx, p)

	expectedPokemon := testPokemon()
	d := Database{
		pool: p,
		log:  slog.Default(),
	}
	_, err := d.pool.Exec(ctx,
		`INSERT INTO pokedex (entry, pokemon) VALUES ($1, $2)`, 1, expectedPokemon,
	)
	assert.NilError(t, err, "failed to insert test pokemon")

	err = d.DeletePokemon(ctx, expectedPokemon.Id)
	assert.NilError(t, err)
}

func TestUpdatePokemon(t *testing.T) {
	ctx := context.Background()
	p := setup(ctx)
	defer cleanup(ctx, p)

	expectedPokemon := testPokemon()
	d := Database{
		pool: p,
		log:  slog.Default(),
	}
	_, err := d.pool.Exec(ctx,
		`INSERT INTO pokedex (entry, pokemon) VALUES ($1, $2)`, 1, expectedPokemon,
	)
	assert.NilError(t, err, "failed to insert test pokemon")

	expectedPokemon.ImageUrl = "shiny_bulbasaur.png"
	got, err := d.UpdatePokemon(ctx, 1, expectedPokemon)
	assert.NilError(t, err)
	assert.DeepEqual(t, got, expectedPokemon)
}
