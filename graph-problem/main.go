package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	uri = "neo4j://localhost:7687"
	// Default username and password are neo4j/neo4j
	// Go to http://localhost:7474 to change the password to "test"
	username = "neo4j"
	password = "12345678"
)

func main() {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(username, password, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection established.")

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	seedData(ctx, session)

	getRecommendations(ctx, session, "user1")

}

func getRecommendations(ctx context.Context, session neo4j.SessionWithContext, userID string) error {
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		records, err := tx.Run(ctx, `
			MATCH (u:User {id: $userID})-[:PURCHASED]->(p:Product)<-[:PURCHASED]-(other:User)-[:PURCHASED]->(rec:Product)
			WHERE NOT (u)-[:PURCHASED]->(rec)
			RETURN rec.id AS recommendedProduct, COUNT(*) AS frequency
			ORDER BY frequency DESC
			LIMIT 10
		`, map[string]interface{}{"userID": userID})

		if err != nil {
			return nil, fmt.Errorf("failed to get recommendations: %w", err)
		}

		list, err := records.Collect(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to collect records: %w", err)
		}

		for _, record := range list {
			fmt.Printf("Recommended product: %+v\n", record.AsMap())
		}

		return nil, records.Err()
	})

	return err
}

func seedData(ctx context.Context, session neo4j.SessionWithContext) {
	// Clean up existing data and seedData the database with some data
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, "MATCH (n) DETACH DELETE n", map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to delete existing data: %w", err)
		}

		_, err = tx.Run(ctx, `
			CREATE (:User {id: 'user1'});
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create nodes: %w", err)
		}

		_, err = tx.Run(ctx, `
			CREATE (:User {id: 'user2'});
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create nodes: %w", err)
		}

		_, err = tx.Run(ctx, `
			CREATE (:Product {id: 'productA'});
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create nodes: %w", err)
		}

		_, err = tx.Run(ctx, `
			CREATE (:Product {id: 'productB'});
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create nodes: %w", err)
		}

		_, err = tx.Run(ctx, `
			CREATE (:Product {id: 'productC'});
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create nodes: %w", err)
		}

		_, err = tx.Run(ctx, `
			CREATE (:Product {id: 'productD'});
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create nodes: %w", err)
		}

		_, err = tx.Run(ctx, `
			MATCH (u:User {id: 'user1'}), (p:Product {id: 'productA'}) CREATE (u)-[:PURCHASED]->(p);
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create relationships: %w", err)
		}

		_, err = tx.Run(ctx, `
			MATCH (u:User {id: 'user1'}), (p:Product {id: 'productB'}) CREATE (u)-[:PURCHASED]->(p);
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create relationships: %w", err)
		}

		_, err = tx.Run(ctx, `
			MATCH (u:User {id: 'user2'}), (p:Product {id: 'productA'}) CREATE (u)-[:PURCHASED]->(p);
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create relationships: %w", err)
		}

		_, err = tx.Run(ctx, `
			MATCH (u:User {id: 'user2'}), (p:Product {id: 'productC'}) CREATE (u)-[:PURCHASED]->(p);
		`, map[string]interface{}{})

		if err != nil {
			return nil, fmt.Errorf("failed to create relationships: %w", err)
		}

		return nil, nil
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Database seeded.")
}
