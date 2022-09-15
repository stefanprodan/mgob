package backup

import (
	"testing"

	"github.com/stefanprodan/mgob/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_BuildRestoreCmd_Different_Name(t *testing.T) {
	target := config.Target{
		Host:     "localhost",
		Port:     27017,
		Username: "user",
		Database: "test",
		Password: "password",
		Params:   "--authenticationDatabase admin",
	}
	restore := config.Target{
		Host:     "localhost",
		Port:     27017,
		Database: "test-restore",
	}

	dumpCmd := BuildRestoreCmd("test.gz", target, restore)
	assert.Equal(t, dumpCmd, `mongorestore --archive=test.gz --gzip --host localhost --port 27017 --nsInclude test.* --nsFrom test.* --nsTo test-restore.* `)
}

func Test_BuildRestoreCmd_Same_Name(t *testing.T) {
	target := config.Target{
		Host:     "localhost",
		Port:     27017,
		Username: "user",
		Database: "test",
		Password: "password",
		Params:   "--authenticationDatabase admin",
	}
	restore := config.Target{
		Host:     "localhost",
		Port:     27017,
		Database: "test",
	}

	dumpCmd := BuildRestoreCmd("test.gz", target, restore)
	assert.Equal(t, dumpCmd, `mongorestore --archive=test.gz --gzip --host localhost --port 27017 --nsInclude test.* `)
}

func Test_BuildDumpCmd_no_uri(t *testing.T) {
	target := config.Target{
		Host:     "localhost",
		Port:     27017,
		Username: "user",
		Database: "test",
		Password: "password",
		Params:   "--authenticationDatabase admin",
	}

	dumpCmd := BuildDumpCmd("test.gz", target)
	assert.Equal(t, dumpCmd, `mongodump --archive=test.gz --gzip --host localhost --port 27017 -u "user" -p "password" --db test --authenticationDatabase admin `)
}

func Test_BuildDumpCmd_uri(t *testing.T) {
	target := config.Target{
		Database: "test",
		Uri:      "mongodb://user:password@localhost:27017",
		Params:   "--authenticationDatabase admin",
	}

	dumpCmd := BuildDumpCmd("test.gz", target)
	assert.Equal(t, dumpCmd, `mongodump --archive=test.gz --gzip --uri "mongodb://user:password@localhost:27017" --db test --authenticationDatabase admin `)
}

func Test_BuildUri_Username(t *testing.T) {
	target := config.Target{
		Host:     "localhost",
		Port:     27017,
		Username: "user",
		Database: "test",
		Password: "password",
	}

	uri := BuildUri(target)
	assert.Equal(t, uri, `mongodb://user:password@localhost:27017`)
}

func Test_BuildUri(t *testing.T) {
	target := config.Target{
		Host:     "localhost",
		Port:     27017,
		Database: "test",
	}

	uri := BuildUri(target)
	assert.Equal(t, uri, `mongodb://localhost:27017`)
}
