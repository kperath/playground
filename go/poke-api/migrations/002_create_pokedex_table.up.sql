-- 002_create_pokedex_table.sql
--
CREATE TABLE pokedex (
    id SERIAL PRIMARY KEY,
    entry INT UNIQUE NOT NULL,
    pokemon JSONB NOT NULL
);

INSERT INTO pokedex (entry, pokemon)
VALUES (
    1,
    '{"pokedex_number":1,"name":"Bulbasaur","type1":"Grass","type2":"Poison","abilities":["Overgrow","Chlorophyll"],"image":"bulbasaur.png"}'
)
