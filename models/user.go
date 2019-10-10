package models

import (
	"database/sql"
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
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

func (o *User) Login(con bolt.Conn) error {
	// Create query
	cypher := `
		  MATCH (n:Comercial)
		  WHERE n.username = {username}
		  RETURN n`

	// Create query params
	params := map[string]interface{}{"username": o.Username}

	st := utils.PrepareSatement(cypher, con)
	rows := utils.QueryStatement(st, params)
	defer func() {
		_ = st.Close()
	}()

	// check results
	row, _, _ := rows.NextNeo()
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

func (o *User) GetNumberOfUsers(con bolt.Conn) (int64, error) {
	var number int64

	// Create query
	cypher := `MATCH (n) RETURN count(*)`
	st := utils.PrepareSatement(cypher, con)
	rows := utils.QueryStatement(st, nil)
	defer func() {
		_ = st.Close()
	}()

	// check results
	row, _, _ := rows.NextNeo()
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

func (o *User) GetNumberOfInvitedUsers(con bolt.Conn) (int64, error) {
	var number int64

	// Create query
	cypher := `MATCH (n)-[r:INVITE]->() WHERE id(n) = {id} RETURN COUNT(r)`

	// Create query params
	params := map[string]interface{}{"id": o.ID}

	st := utils.PrepareSatement(cypher, con)
	rows := utils.QueryStatement(st, params)
	defer func() {
		_ = st.Close()
	}()

	// check results
	row, _, _ := rows.NextNeo()
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

	con := utils.CreateConnection()

	// Create query
	cypher := `
	  MATCH (n:Comercial)
      WHERE n.username = {username}
	  RETURN n`

	// Prepare query
	params := map[string]interface{}{"username": o.Username}

	st := utils.PrepareSatement(cypher, con)
	rows := utils.QueryStatement(st, params)
	defer func() {
		_ = st.Close()
	}()

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
		con := utils.CreateConnection()

		log.Println("no")

		cypher := `CREATE (m:Comercial {username: {username}, password: {password}, name: {name}, api_token: "" })`
		params := map[string]interface{}{
			"name":     o.Name,
			"username": o.Username,
			"password": utils.HashAndSalt([]byte(o.Password))}

		st := utils.PrepareSatement(cypher, con)
		utils.ExecuteStatement(st, params)
		defer func() {
			_ = st.Close()
		}()
	}

	return nil
}

func (o *User) UpdateToken(apiToken string) error {

	con, err := bolt.NewDriver().OpenNeo("bolt://neo4j:123456@localhost:7687")
	utils.HandleError(err)
	defer func() {
		_ = con.Close()
	}()

	// Create statement
	cypher := `MATCH (n:Comercial { username: {username} })
				SET n.api_token = {api_token}
				RETURN n.api_token`

	// Create params
	params := map[string]interface{}{"username": o.Username, "api_token": apiToken}

	st := utils.PrepareSatement(cypher, con)
	utils.ExecuteStatement(st, params)
	defer func() {
		_ = st.Close()
	}()

	return nil
}

func (o *User) InviteUser(con bolt.Conn, invitedID int64) error {

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

	st := utils.PrepareSatement(cypher, con)
	utils.ExecuteStatement(st, params)
	defer func() {
		_ = st.Close()
	}()

	return nil
}

func (o *User) GetUserByID(con bolt.Conn) error {

	// Create query
	cypher := `MATCH (n:Comercial) where id(n) = {id} RETURN n`

	// Create query params
	params := map[string]interface{}{"id": o.ID}

	st := utils.PrepareSatement(cypher, con)
	rows := utils.QueryStatement(st, params)
	defer func() {
		_ = st.Close()
	}()
	// check results
	row, _, _ := rows.NextNeo()
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

func (o *User) GetAll(con bolt.Conn) ([]User, error) {
	var list = make([]User, 0)

	// Create query
	cypher := `
	  MATCH (n:Comercial)
	  RETURN n`

	st := utils.PrepareSatement(cypher, con)
	rows := utils.QueryStatement(st, nil)
	defer func() {
		_ = st.Close()
	}()

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

func (o *User) GetNumberOfUsers2(con bolt.Conn) (int64, error) {
	var number int64

	// Create query
	cypher := `MATCH (n) RETURN count(*)`

	st := utils.PrepareSatement(cypher, con)
	rows := utils.QueryStatement(st, nil)
	defer func() {
		_ = st.Close()
	}()

	// check results
	row, _, _ := rows.NextNeo()
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
