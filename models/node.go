package models

type Neo struct {
	NodeIdentity int      `json:"NodeIdentity"`
	Labels       []string `json:"Labels"`
	Properties   struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"Properties"`
}