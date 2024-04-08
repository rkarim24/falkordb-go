package falkordb

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type FalkorDB struct {
	Conn *redis.Client
}

type ConnectionOption = redis.Options

func isSentinel(conn *redis.Client) bool {
	info, _ := conn.InfoMap(ctx, "server").Result()
	return info["Server"]["redis_mode"] == "sentinel"
}

// FalkorDB Class for interacting with a FalkorDB server.
func FalkorDBNew(options *ConnectionOption) (*FalkorDB, error) {
	db := redis.NewClient(options)

	if isSentinel(db) {
		masters, err := db.Do(ctx, "SENTINEL", "MASTERS").Result()
		if err != nil {
			return nil, err
		}
		if len(masters.([]interface{})) != 1 {
			return nil, errors.New("multiple masters, require service name")
		}
		str := "name"
		var strInterface interface{} = str
		masterName := masters.([]interface{})[0].(map[interface{}]interface{})[strInterface].(string)
		db = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    masterName,
			SentinelAddrs: []string{options.Addr},
		})
	}
	return &FalkorDB{Conn: db}, nil
}

// Creates a new FalkorDB instance from a URL.
func FromURL(url string) (*FalkorDB, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	db := redis.NewClient(options)
	if isSentinel(db) {
		masters, err := db.Do(ctx, "SENTINEL", "MASTERS").Result()
		if err != nil {
			return nil, err
		}
		if len(masters.([]interface{})) != 1 {
			return nil, errors.New("multiple masters, require service name")
		}
		masterName := masters.([]interface{})[0].(map[string]interface{})["name"].(string)
		db = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    masterName,
			SentinelAddrs: []string{options.Addr},
		})
	}
	return &FalkorDB{Conn: db}, nil
}

// Selects a graph by creating a new Graph instance.
func (db *FalkorDB) SelectGraph(graphName string) *Graph {
	return graphNew(graphName, db.Conn)
}

// List all graph names.
// See: https://docs.falkordb.com/commands/graph.list.html
func (db *FalkorDB) ListGraphs() ([]string, error) {
	return db.Conn.Do(ctx, "GRAPH.LIST").StringSlice()
}

// Retrieve a DB level configuration.
// For a list of available configurations see: https://docs.falkordb.com/configuration.html#falkordb-configuration-parameters
func (db *FalkorDB) ConfigGet(key string) string {
	return db.Conn.Do(ctx, "GRAPH.CONFIG", "GET", key).String()
}

// Update a DB level configuration.
// For a list of available configurations see: https://docs.falkordb.com/configuration.html#falkordb-configuration-parameters
func (db *FalkorDB) ConfigSet(key, value string) error {
	return db.Conn.Do(ctx, "GRAPH.CONFIG", "SET", key).Err()
}
