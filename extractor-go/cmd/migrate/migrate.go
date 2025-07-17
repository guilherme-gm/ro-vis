/**
 * "Migrate" entry point
 * Applies DB migrations
 */

package main

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/guilherme-gm/ro-vis/extractor/internal/conf"
	"github.com/guilherme-gm/ro-vis/extractor/internal/domain/server"
)

type MigrationTool struct {
	tool   *migrate.Migrate
	server *server.Server
}

func newMigrationTool(server *server.Server) (*MigrationTool, error) {

	m, err := migrate.New(
		"file://../migrations",
		strings.Replace(conf.MigratorConfig.DbUrl, "@DATABASE@", server.DatabaseName, 1))

	if err != nil {
		return nil, err
	}

	return &MigrationTool{tool: m, server: server}, nil
}

func (m *MigrationTool) UpdateCheck() {
	fmt.Println("-- Checking for updates for " + m.server.DatabaseName + " --")
	version, dirty, err := m.tool.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		fmt.Println("Empty database.")
		return
	}

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
	fmt.Println("-- Updating " + m.server.DatabaseName + " --")
	err := m.tool.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Database is up to date.")
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Println("Database updated successfully.")

	return nil
}

func main() {
	fmt.Println("RO Vis extractor - Migrate")
	conf.LoadMigrator()

	for _, server := range server.GetServers() {
		fmt.Println("------ Migrating " + server.DatabaseName + " ------")
		migTool, err := newMigrationTool(server)
		if err != nil {
			fmt.Printf("Failed to create migration tool for %s: %v\n", server.DatabaseName, err)
			return
		}

		migTool.UpdateCheck()

		err = migTool.Up()
		if err != nil {
			fmt.Printf("Failed to update %s: %v\n", server.DatabaseName, err)
			return
		}
	}

	fmt.Println("Success")
}
