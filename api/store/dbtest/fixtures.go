package dbtest

import "github.com/shellhub-io/mongotest"

const FixtureUser = "user"

func (db *DB) ApplyFixtures(fixtures ...string) error {
	return mongotest.UseFixture(fixtures...)
}

func (db *DB) TeardownFixtures(fixtures ...string) error {
	return mongotest.DropDatabase()
}
