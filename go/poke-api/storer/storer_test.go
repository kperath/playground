package storer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"

	"playground/poke-api/types"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"gotest.tools/v3/assert"
)

func setupDB(t *testing.T, ctx context.Context) *pgxpool.Pool {
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
		t.Fatalf("creating table: %s", err)
	}
	return pool
}

func cleanupDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	defer pool.Close()
	_, err := pool.Exec(ctx, `DROP TABLE IF EXISTS pokedex;`)
	if err != nil {
		t.Fatal(err)
	}
}

func testPokemon() *types.Pokemon {
	return &types.Pokemon{
		Id:        1,
		Name:      "Bulbasaur",
		Type1:     "Grass",
		Type2:     "Poison",
		Abilities: []string{"Overgrow", "Chlorophyll"},
		ImageUrl:  "bulbasaur.png",
	}
}

func TestGetPokemon(t *testing.T) {
	ctx := context.Background()
	p := setupDB(t, ctx)
	defer cleanupDB(t, ctx, p)

	expectedPokemon := testPokemon()
	d := &Database{
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

func TestAddPokemonPostgres(t *testing.T) {
	ctx := context.Background()
	p := setupDB(t, ctx)
	defer cleanupDB(t, ctx, p)

	expectedPokemon := testPokemon()
	d := &Database{
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
	p := setupDB(t, ctx)
	defer cleanupDB(t, ctx, p)

	expectedPokemon := testPokemon()
	d := &Database{
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
	p := setupDB(t, ctx)
	defer cleanupDB(t, ctx, p)

	expectedPokemon := testPokemon()
	d := &Database{
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

func setupElastic(t *testing.T) *elasticsearch.TypedClient {
	ctx := context.Background()
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("TEST_ELASTIC_URL")},
		Username:  "elastic",
		Password:  "elastic",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Indices.Get("pokedex").Do(ctx)
	_, ok := resp["pokedex"]
	if ok && err != nil {
		t.Fatal(err)
	}
	if !ok {
		// Create the index with settings
		f, err := os.Open("../pokeapi-index.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		_, err = client.Indices.Create("pokedex").Raw(f).Do(ctx)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Indices.Refresh().Index("pokedex").Do(ctx)
		if err != nil {
			t.Fatalf("failed to refresh index: %v", err)
		}
	}
	return client
}

func cleanupElastic(t *testing.T, ctx context.Context, client *elasticsearch.TypedClient) {
	_, err := client.Indices.Delete("pokedex").Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddPokemonElastic(t *testing.T) {
	ctx := context.Background()
	client := setupElastic(t)
	defer cleanupElastic(t, ctx, client)

	s := &Search{
		client: client,
		log:    slog.Default(),
	}
	expectedPokemon := testPokemon()
	err := s.AddPokemon(ctx, expectedPokemon)
	assert.NilError(t, err)

	resp, err := client.
		Get("pokedex", fmt.Sprintf("%d", expectedPokemon.Id)).
		Do(ctx)
	assert.NilError(t, err)

	gotPokemon := &types.Pokemon{}
	err = json.Unmarshal(resp.Source_, &gotPokemon)
	assert.DeepEqual(t, gotPokemon, expectedPokemon)
}

func TestSearchPokemon(t *testing.T) {
	ctx := context.Background()
	client := setupElastic(t)
	defer cleanupElastic(t, ctx, client)

	s := &Search{
		client: client,
		log:    slog.Default(),
	}

	testPokemon := testPokemon()
	_, err := client.Index("pokedex").
		Id(fmt.Sprintf("%d", testPokemon.Id)).
		Document(testPokemon).
		Do(context.Background())
	assert.NilError(t, err)

	// refresh index so that data is immediately queriable
	_, err = client.Indices.Refresh().Index("pokedex").Do(context.Background())
	assert.NilError(t, err)

	type test struct {
		name string
	}
	for _, tc := range []test{
		{
			name: "bulba",
		},
		{
			name: "bulbasaur",
		},
		{
			name: "Bulbasaur",
		},
		{
			name: "Bulba",
		},
		{
			name: "B",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// NOTE: don't parallelize tests due to cleaning up of index
			p, resultCount, err := s.SearchPokemon(ctx, tc.name, 1, 1)
			assert.Equal(t, resultCount, int64(1))
			assert.NilError(t, err)
			assert.Equal(t, len(p), 1)
			assert.DeepEqual(t, p[0], testPokemon)
		})
	}
}
