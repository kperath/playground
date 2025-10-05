# Poke API

My example API that I will use in all my projects to test things, maybe it should be its own repo, maybe not. If I finish it, then it will be.

So basically I need to create a CRUD app!

GET /pokedex/search?query=&page=

GET /pokedex

DELETE /admin/pokedex
{
    id: 1
}

POST /admin/pokedex
{
    id: 1
    new: json data
}

PATCH /admin/pokedex
{
    id: 1
    change: json data
}

##### Architecture

(Client)
-->  pokedex_service (search queries)       --> ElasticSearch
-->  pokemon_service                        --> DB cdc to ElasticSearch (?)


##### Future

Pokemon Party (limited to 6 entries)

GET /pokemon/party

PUT /pokemon/party
{
new_members: {
    1: {pokemon},
    2: {pokemon},
}
}

POST /pokemon/party
{
party: {
initial pokemon party
    1: {pokemon},
    2: {pokemon}
}
}

DELETE /pokemon/party
{ id: 1 }


##### Architecture

(Client)
-->  pokedex_service (search queries)       --> ElasticSearch
-->  pokemon_service                        --> DB cdc to ElasticSearch (?)
