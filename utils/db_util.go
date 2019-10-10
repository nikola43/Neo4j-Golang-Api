package utils

import (
	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"log"
)

func UpdateToken(username string, api_token string) {
	// open connection
	db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer func() {
		_ = db.Close()
	}()

	cypher :=  `MATCH (n:Comercial { username: {username} })
				SET n.api_token = {api_token}
				RETURN n.name`

	params := map[string]interface{}{"username": username, "api_token": api_token}

	createErr, result := db.ExecNeo(cypher, params)
	if createErr != nil {
		log.Println(createErr)
	}

	log.Println(result)
}

func CheckIfUsersExistsOnDB(username string) bool {
	var exists = false

	// Open connection
	db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer func() {
		_ = db.Close()
	}()

	// Create query
	cypher := `
	  MATCH (n:Comercial)
      WHERE n.username = {username}
	  RETURN n`

	// Prepare query
	params := map[string]interface{}{"username": username}

	// Make query
	rows, err := db.QueryNeo(cypher, params)
	if err != nil {
		log.Println("error querying graph:", err)
	}

	// check if has result
	row, _, err := rows.NextNeo()
	if len(row) > 0 {
		exists = true
	}

	return exists
}
