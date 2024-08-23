package test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDatabaseAndPool(connString string, database string) (*pgxpool.Pool, error) {
	c, err := PostgresConn(connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres database: %w", err)
	}
	defer c.Close(context.Background())

	if err := CreateDatabase(c, database); err != nil {
		return nil, fmt.Errorf("error creating conn to %s database: %w", database, err)

	}

	p, err := Pool(connString, database)
	if err != nil {
		return nil, fmt.Errorf("error creating pool to %s database: %w", database, err)
	}

	if err := MigrateDatabase(p); err != nil {
		p.Close()
		return nil, fmt.Errorf("error migrating database %s : %w", database, err)
	}

	return p, nil
}

func TeardownDatabase(connString string, database string) error {
	c, err := PostgresConn(connString)
	if err != nil {
		return err
	}
	defer c.Close(context.Background())

	err = DropDatabase(c, database)
	if err != nil {
		return err
	}
	return nil
}
