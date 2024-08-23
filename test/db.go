package test

import (
	"accountant"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Pool(connString string, database string) (*pgxpool.Pool, error) {
	c, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing db config: %w", err)
	}
	c.ConnConfig.TLSConfig = nil
	c.ConnConfig.Database = database

	p, err := pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("error creating db pool: %w", err)
	}

	return p, nil
}

func PostgresConn(connString string) (*pgx.Conn, error) {
	conf, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing db config: %w", err)
	}
	conf.TLSConfig = nil
	conf.Database = "postgres"

	conn, err := pgx.ConnectConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return conn, nil
}

func CreateDatabase(conn *pgx.Conn, name string) error {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (force)", name))
	if err != nil {
		return fmt.Errorf("error during create database step=drop: %w", err)
	}
	_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", name))
	if err != nil {
		return fmt.Errorf("error during create database step=create: %w", err)
	}

	return nil
}

func DropDatabase(conn *pgx.Conn, name string) error {
	_, err := conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (force)", name))
	if err != nil {
		return fmt.Errorf("error during drop database: %w", err)
	}

	return nil
}

func DeleteFromTables(pool *pgxpool.Pool, tables ...string) error {
	for _, table := range tables {
		_, err := pool.Exec(context.Background(), fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("error during deleting from %s: %w", table, err)
		}
	}

	return nil
}

func MigrateDatabase(pool *pgxpool.Pool) error {
	goose.SetBaseFS(accountant.EmbedFs)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("error during setting goose dialect: %w", err)
	}

	if err := goose.Up(stdlib.OpenDBFromPool(pool), "migrations"); err != nil {
		return fmt.Errorf("error during migrating db: %w", err)
	}

	return nil
}
