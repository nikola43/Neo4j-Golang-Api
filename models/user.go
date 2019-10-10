package models

import (
	"database/sql"
	"fmt"
	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
	"github.com/nikola43/ecodadys_api/utils"
	"log"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	ApiToken string `json:"api_token"`
}

func (o *User) Login() error {

	db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return err
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
	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Println("error preparing graph:", err)
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	// Create query params
	query := map[string]interface{}{"username": o.Username}

	// Make query
	rows, err := stmt.QueryNeo(query)
	if err != nil {
		log.Println("error querying graph:", err)
		return err
	}

	// check results
	row, _, err := rows.NextNeo()
	if len(row) > 0 {
		for _, element := range row {
			currentNode := element.(graph.Node)
			username := currentNode.Properties["username"].(string)
			password := currentNode.Properties["password"].(string)

			// check if user and password match
			if username == o.Username && utils.ComparePasswords(password, []byte(o.Password)) {
				apiToken := utils.GenerateTokenUsername(o.Username)
				err := o.UpdateToken(apiToken)
				if err != nil {
					log.Println("error updating user token", err)
					return err
				}

				// update user info
				o.ID = currentNode.NodeIdentity
				o.Name = currentNode.Properties["name"].(string)
				o.Password = ""
				o.ApiToken = apiToken

				log.Print("encontrado: ")
				log.Println(o.ID)
			} else {
				log.Println("credentials not match")
				err := sql.ErrNoRows
				return err
			}
		}
	} else {
		log.Println("not found")
		err := fmt.Errorf("not foundd")
		return err
	}

	return nil
}

func (o *User) GetNumberOfUsers() (int64, error) {
	var number int64

	db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return 0, err
	}
	defer func() {
		_ = db.Close()
	}()

	// Create query
	cypher := `MATCH (n) RETURN count(*)`

	// Prepare query
	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Println("error preparing graph:", err)
		return 0, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	// Make query
	rows, err := stmt.QueryNeo(nil)
	if err != nil {
		log.Println("error querying graph:", err)
		return 0, err
	}

	// check results
	row, _, err := rows.NextNeo()
	if len(row) > 0 {
		log.Println("sss found")
		for _, element := range row {
			fmt.Println(row)
			number = element.(int64)
		}
	} else {
		log.Println("not found")
		err := sql.ErrNoRows
		return 0, err
	}

	return number, nil
}

func (o *User) GetNumberOfInvitedUsers() (int64, error) {
	var number int64

	db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return 0, err
	}
	defer func() {
		_ = db.Close()
	}()

	// Create query
	cypher := `MATCH (n)-[r:INVITE]->() WHERE id(n) = {id} RETURN COUNT(r)`

	// Prepare query
	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Println("error preparing graph:", err)
		return 0, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	// Create query params
	query := map[string]interface{}{"id": o.ID}

	// Make query
	rows, err := stmt.QueryNeo(query)
	if err != nil {
		log.Println("error querying graph:", err)
		return 0, err
	}

	// check results
	row, _, err := rows.NextNeo()
	if len(row) > 0 {
		for _, element := range row {
			number = element.(int64)
		}
	} else {
		log.Println("not found")
		err := sql.ErrNoRows
		return 0, err
	}

	return number, nil
}

func (o *User) SingUp() error {
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
	params := map[string]interface{}{"username": o.Username}

	// Make query
	rows, err := db.QueryNeo(cypher, params)
	if err != nil {
		log.Println("error querying graph:", err)
	}

	// Check Results
	row, _, err := rows.NextNeo()
	if len(row) > 0 {
		log.Println("Hay datos")
		for _, element := range row {
			currentNode := element.(graph.Node)
			if currentNode.Properties["username"].(string) == o.Username {
				log.Println("Ya existe")
				err = fmt.Errorf("username %s already exists", o.Username)
				return err
			} else {
				log.Println("no")
			}
			row, _, err = rows.NextNeo()
		}
	} else {
		// open connection
		db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
		if err != nil {
			log.Println("error connecting to neo4j:", err)
		}
		defer func() {
			_ = db.Close()
		}()

		log.Println("no")

		hashPassword := utils.HashAndSalt([]byte(o.Password))

		cypher := `CREATE (m:Comercial {username: {username}, password: {password}, name: {name}, api_token: "" })`
		params := map[string]interface{}{
			"name":     o.Name,
			"username": o.Username,
			"password": hashPassword}

		fmt.Println(hashPassword)

		result, createErr := db.ExecNeo(cypher, params)
		if createErr != nil {
			log.Println(createErr)
		}

		numResult, err := result.RowsAffected()
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println(numResult)
	}

	return nil
}

func (o *User) UpdateToken(apiToken string) error {
	// open connection
	db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	// Create statement
	cypher := `MATCH (n:Comercial { username: {username} })
				SET n.api_token = {api_token}
				RETURN n.api_token`

	// Create params
	params := map[string]interface{}{"username": o.Username, "api_token": apiToken}

	// Exec statement
	result, updateTokenError := db.ExecNeo(cypher, params)
	if updateTokenError != nil {
		log.Println(updateTokenError)
		return updateTokenError
	}

	numResult, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(numResult)
	return nil
}

func (o *User) InviteUser(invitedID int64) error {
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
	  MATCH (a:Comercial),(b:Comercial)
      WHERE id(a) = {invite_id} AND id(b) = {invited_id} 
      CREATE (a)-[r:INVITE]->(b)
      RETURN r`

	// Prepare query
	params := map[string]interface{}{
		"invite_id":  o.ID,
		"invited_id": invitedID}

	result, createErr := db.ExecNeo(cypher, params)
	if createErr != nil {
		log.Println(createErr)
	}

	numResult, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(numResult)
	if numResult < 1 {
		err := fmt.Errorf("error creating relationship")
		return err
	}

	return nil
}

func (o *User) GetUserByID() error {

	db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	// Create query
	cypher := `MATCH (n:Comercial) where id(n) = {id} RETURN n`

	// Prepare query
	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Println("error preparing graph:", err)
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	// Create query params
	query := map[string]interface{}{"id": o.ID}

	// Make query
	rows, err := stmt.QueryNeo(query)
	if err != nil {
		log.Println("error querying graph:", err)
		return err
	}

	// check results
	row, _, err := rows.NextNeo()
	if len(row) > 0 {
		for _, element := range row {
			currentNode := element.(graph.Node)
			o.ID = currentNode.NodeIdentity
			o.Name = currentNode.Properties["name"].(string)
			o.ApiToken = currentNode.Properties["api_token"].(string)
		}
	} else {
		log.Println("not found")
		err := sql.ErrNoRows
		return err
	}

	return nil
}

func (o *User) GetAll() ([]User, error) {
	var list = make([]User, 0)

	db, err := driver.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return list, err
	}
	defer func() {
		_ = db.Close()
	}()

	// Create query
	cypher := `
	  MATCH (n:Comercial)
	  RETURN n`

	// Prepare query
	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Println("error preparing graph:", err)
		return list, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	// Create query params
	query := map[string]interface{}{}

	// Make query
	rows, err := stmt.QueryNeo(query)
	if err != nil {
		log.Println("error querying graph:", err)
		return list, err
	}

	// check results
	row, _, err := rows.NextNeo()

	for row != nil && err == nil {
		for _, element := range row {
			currentNode := element.(graph.Node)
			var u User
			u.ID = currentNode.NodeIdentity
			u.Name = currentNode.Properties["name"].(string)
			u.Username = currentNode.Properties["username"].(string)
			u.Password = currentNode.Properties["password"].(string)
			u.ApiToken = currentNode.Properties["api_token"].(string)
			list = append(list, u)
			fmt.Println(u)
			row, _, err = rows.NextNeo()
		}
	}

	for _, element := range row {
		currentNode := element.(graph.Node)
		var u User
		u.ID = currentNode.NodeIdentity
		u.Name = currentNode.Properties["name"].(string)
		u.Username = currentNode.Properties["username"].(string)
		u.Password = currentNode.Properties["password"].(string)
		u.ApiToken = currentNode.Properties["api_token"].(string)
		list = append(list, u)
		row, _, err = rows.NextNeo()
	}

	return list, nil
}
