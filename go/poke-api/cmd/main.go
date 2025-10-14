package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"playground/poke-api/api/pokedex"
	"playground/poke-api/storer"
	"playground/poke-api/types"
	"strconv"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
)

func main() {
	log := log.New(os.Stdout, "setup:", 0)

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic("connecting to db", err)
	}
	defer pool.Close()

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTIC_URL")},
	})

	ll := slog.Default()

	es := storer.NewSearch(ll, client)
	db := storer.NewDatabase(ll, pool)

	dex := pokedex.New(db, es)

	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/pokedex/:id", func(c echo.Context) error {
		i := c.Param("id")
		id, err := strconv.Atoi(i)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, types.Pokemon{})
		}

		p, err := dex.GetPokemon(c.Request().Context(), id)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, types.Pokemon{})
		}
		return c.JSON(http.StatusOK, p)
	})

	e.GET("/pokedex/search", func(c echo.Context) error {
		nameQuery := c.QueryParam("name")
		pageQuery := c.QueryParam("page")
		pageSizeQuery := c.QueryParam("page_size")

		var pageOpts *pokedex.SearchOpts
		if pageQuery != "" || pageSizeQuery != "" {
			page, err := strconv.Atoi(pageQuery)
			if err != nil {
				e.Logger.Error(err)
				return c.JSON(http.StatusNotFound, pokedex.PaginatedPokemon{})
			}
			pageSize, err := strconv.Atoi(pageSizeQuery)
			if err != nil {
				e.Logger.Error(err)
				return c.JSON(http.StatusNotFound, pokedex.PaginatedPokemon{})
			}
			pageOpts = &pokedex.SearchOpts{
				Page:     page,
				PageSize: pageSize,
			}
		}
		p, err := dex.SearchPokemon(c.Request().Context(), nameQuery, pageOpts)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, pokedex.PaginatedPokemon{})
		}
		return c.JSON(http.StatusOK, p)
	})

	e.POST("/admin/pokedex", func(c echo.Context) error {
		req := c.Request()
		defer req.Body.Close()

		p := &types.Pokemon{}
		err := json.NewDecoder(req.Body).Decode(p)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		err = dex.AddPokemon(c.Request().Context(), p)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		return c.JSON(http.StatusAccepted, nil)
	})

	e.PUT("/admin/pokedex/:id", func(c echo.Context) error {
		i := c.Param("id")
		id, err := strconv.Atoi(i)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}
		req := c.Request()
		defer req.Body.Close()

		var p *types.Pokemon
		err = json.NewDecoder(req.Body).Decode(p)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		p, err = dex.UpdatePokemon(c.Request().Context(), id, p)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusAccepted, p)
	})

	e.DELETE("/admin/pokemon/:id", func(c echo.Context) error {
		i := c.Param("id")
		id, err := strconv.Atoi(i)
		if err != nil {
			e.Logger.Error(err)
			return c.JSON(http.StatusBadRequest, nil)
		}

		err = dex.DeletePokemon(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusAccepted, nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
	// test handlers
	// TODO: add in auth middleware that guards admin endpoints
	// TODO: /login & JWT
	// TODO: prometheus instrumentation
	// TODO extra: containerize
	// TODO extra: k8s kind - can stop here tbh
	// kube-prometheus
	// TODO extra: istio
	// TODO extra: cilium
	// TODO extra: cert-manager
	// TODO extra: flux?
}
