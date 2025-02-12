[![license](https://img.shields.io/github/license/FalkorDB/falkordb-go.svg)](https://github.com/FalkorDB/falkordb-go)
[![GitHub issues](https://img.shields.io/github/release/FalkorDB/falkordb-go.svg)](https://github.com/FalkorDB/falkordb-go/releases/latest)
[![Codecov](https://codecov.io/gh/FalkorDB/falkordb-go/branch/master/graph/badge.svg)](https://codecov.io/gh/FalkorDB/falkordb-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/FalkorDB/falkordb-go)](https://goreportcard.com/report/github.com/FalkorDB/falkordb-go)
[![GoDoc](https://godoc.org/github.com/FalkorDB/falkordb-go?status.svg)](https://godoc.org/github.com/FalkorDB/falkordb-go)

# falkordb-go
[![Discord](https://img.shields.io/discord/1146782921294884966?style=flat-square)](https://discord.gg/6M4QwDXn2w)

`falkordb-go` is a Golang client for the [FalkorDB](https://falkordb.com) database.

## Installation

Simply do:
```sh
$ go get github.com/FalkorDB/falkordb-go
```

## Usage

The complete `falkordb-go` API is documented on [GoDoc](https://godoc.org/github.com/FalkorDB/falkordb-go).

```go
package main

import (
	"fmt"
	"os"

	"github.com/FalkorDB/falkordb-go"
)

func main() {
	db, err := falkordb.FalkorDBNew(&falkordb.ConnectionOption{Addr: "0.0.0.0:6379"})
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}

	graph := db.SelectGraph("social")

	_, err = graph.Query("CREATE (:Person {name: 'John Doe', age: 33, gender: 'male', status: 'single'})-[:VISITED]->(:VISITED {name: 'Japan'})", nil, nil)
	if err != nil {
		fmt.Println("Error executing CREATE query:", err)
		os.Exit(1)
	}

	query := "MATCH (p:Person)-[v:VISITED]->(c:VISITED) RETURN p.name, p.age, c.name"
	result, err := graph.Query(query, nil, nil)
	if err != nil {
		fmt.Println("Error executing MATCH query:", err)
		os.Exit(1)
	}

	// Check if result is nil
	if result == nil {
		fmt.Println("No result returned from the query.")
		os.Exit(1)
	}

	// Pretty-print the full result set as a table.
	result.PrettyPrint()

	// Iterate over each individual Record in the result.
	fmt.Println("Visited countries by person:")
	for result.Next() { // Next returns true until the iterator is depleted.
		// Get the current Record.
		r := result.Record()

		// Entries in the Record can be accessed by index or key.
		pName := r.GetByIndex(0)
		fmt.Printf("\nName: %s\n", pName)
		pAge, _ := r.Get("p.age")
		fmt.Printf("\nAge: %d\n", pAge)
	}

	// Path matching example.
	query = "MATCH p = (:Person)-[:VISITED]->(:VISITED) RETURN p"
	result, err = graph.Query(query, nil, nil)
	if err != nil {
		fmt.Println("Error executing MATCH path query:", err)
		os.Exit(1)
	}

	fmt.Println("Paths of persons visiting countries:")
	for result.Next() {
		r := result.Record()
		p, ok := r.GetByIndex(0).(falkordb.Path)
		fmt.Printf("%v %v\n", p, ok)
	}
}
```

Running the above produces the output:

```sh
+----------+-------+--------+
|  p.name  | p.age | c.name |
+----------+-------+--------+
| John Doe |    33 | Japan  |
+----------+-------+--------+

Query internal execution time 1.623063

Name: John Doe

Age: 33
```

## Running queries with timeouts

Queries can be run with a millisecond-level timeout as described in [the documentation](https://docs.falkordb.com/configuration.html#timeout). To take advantage of this feature, the `QueryOptions` struct should be used:

```go
options := NewQueryOptions().SetTimeout(10) // 10-millisecond timeout
res, err := graph.Query("MATCH (src {name: 'John Doe'})-[*]->(dest) RETURN dest", nil, options)
```

## Running tests

A simple test suite is provided, and can be run with:

```sh
$ go test
```

The tests expect a FalkorDB server to be available at localhost:6379

## License

falkordb-go is distributed under the BSD3 license - see [LICENSE](LICENSE)
