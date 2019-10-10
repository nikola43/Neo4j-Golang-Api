package utils

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"

	"log"
)

func CreateConnection() bolt.Conn {
	con, err := bolt.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	HandleError(err)
	return con
}

// Here we prepare a new statement. This gives us the flexibility to
// cancel that statement without any request sent to Neo
func PrepareSatement(query string, con bolt.Conn) bolt.Stmt {
	st, err := con.PrepareNeo(query)
	HandleError(err)
	return st
}

// Executing a statement just returns summary information
func ExecuteStatement(st bolt.Stmt, params map[string]interface{}) {
	result, err := st.ExecNeo(params)
	HandleError(err)
	numResult, err := result.RowsAffected()
	HandleError(err)
	log.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1

	// Closing the statement will also close the rows
	defer func() {
		_ = st.Close()
	}()
}

func QueryStatement(st bolt.Stmt, params map[string]interface{}) bolt.Rows {
	// Even once I get the rows, if I do not consume them and close the
	// rows, Neo will discard and not send the data
	rows, err := st.QueryNeo(params)
	HandleError(err)
	return rows
}

func UpdateToken(username string, api_token string) {
	// open connection
	db, err := bolt.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer func() {
		_ = db.Close()
	}()

	cypher := `MATCH (n:Comercial { username: {username} })
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
	db, err := bolt.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
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
