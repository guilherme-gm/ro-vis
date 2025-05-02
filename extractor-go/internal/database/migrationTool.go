package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
)

type MigrationTool struct {
	tool *migrate.Migrate
}

func NewMigrationTool() (*MigrationTool, error) {
	m, err := migrate.New(
		"file://../migrations",
		"mysql://"+conf.Config.DbUrl)

	if err != nil {
		return nil, err
	}

	return &MigrationTool{tool: m}, nil
}

func (m *MigrationTool) UpdateCheck() {
	version, dirty, err := m.tool.Version()
	if err != nil {
		fmt.Println("------ Database is not ready ------")
		fmt.Println(err)
		panic("Can't continue until above error is fixed.")
	}

	if dirty {
		fmt.Println("------ Database is not ready ------")
		fmt.Println("Database is dirty. You will likely need to perform some manual fix.")
		panic("Can't continue until above error is fixed.")
	}

	fmt.Printf("Running database version: %d\n", version)
	// @TODO: Check to ensure we are in the latest version, we could either scan the filers, or, preferable
	// with this PR once it is merged: https://github.com/golang-migrate/migrate/pull/1219
}

func (m *MigrationTool) Up() error {
	return m.tool.Up()
}
