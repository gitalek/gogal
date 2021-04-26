package main

import "fmt"

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	// We are going to provide two potential connection info strings
	// based on whether a password is present.
	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=dilable",
			c.Host, c.Port, c.User, c.Name,
		)
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=dilable",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "jon",
		Password: "your-password",
		Name:     "gogal_dev",
	}
}
