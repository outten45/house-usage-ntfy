Mix.install([ 
  {:req, "~> 0.3.6"}, 
  {:jason, "~> 1.0"},
  {:exqlite, "~> 0.13.9"}
])

defmodule Test do
end

alias Exqlite.Sqlite3

# Req.get!("https://api.github.com/repos/elixir-lang/elixir").body#["description"]
# |> dbg()

{:ok, conn} = Sqlite3.open("rtl.db")

{:ok, statement} = Exqlite.Sqlite3.prepare(conn, "select id, type,label,value From measurements where label = ?1 limit 10") 
|> IO.inspect(label: ">prepare")
r =  Exqlite.Sqlite3.bind(conn, statement, ["electtric-recv"])
|> IO.inspect(label: ">> r")
dbg(statement)

results = Sqlite3.step(conn, statement)
|> IO.inspect(label: "> results")

dbg(results)
dbg(r)

{:ok, statement} = Exqlite.Sqlite3.prepare(conn, "select id, type,label,value From measurements where label = 'electric-recv' limit 10")
dbg(Exqlite.Sqlite3.step(conn, statement))
dbg(Exqlite.Sqlite3.step(conn, statement))
dbg(Exqlite.Sqlite3.step(conn, statement))

Sqlite3.release(conn, nil)