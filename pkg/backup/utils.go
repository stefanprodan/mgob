package backup

import (
	"fmt"

	"github.com/stefanprodan/mgob/pkg/config"
)

func BuildDumpCmd(archive string, target config.Target) string {
	return buildCmd("mongodump", archive, target, "")
}

func BuildRestoreCmd(archive string, target config.Target, restore config.Target) string {
	if target.Database != restore.Database {
		cmd := buildCmd("mongorestore", archive, restore, target.Database)
		cmd += fmt.Sprintf("--nsFrom %v.* --nsTo %v.* ", target.Database, restore.Database)
		return cmd
	}
	return buildCmd("mongorestore", archive, restore, "")

}

// command: mongodump | mongorestore
func buildCmd(command string, archive string, target config.Target, overrideDBName string) string {
	cmd := fmt.Sprintf("%v --archive=%v --gzip ", command, archive)
	// using uri (New in version 3.4.6)
	// host/port/username/password are incompatible with uri
	// https://docs.mongodb.com/manual/reference/program/mongodump/#cmdoption-mongodump-uri
	// use older host/port
	if target.Uri != "" {
		cmd += fmt.Sprintf(`--uri "%v" `, target.Uri)
	} else {
		cmd += fmt.Sprintf("--host %v --port %v ", target.Host, target.Port)

		if target.Username != "" && target.Password != "" {
			cmd += fmt.Sprintf(`-u "%v" -p "%v" `, target.Username, target.Password)
		}
	}

	if target.Database != "" {
		if command == "mongodump" {
			cmd += fmt.Sprintf("--db %v ", target.Database)
		} else {
			if overrideDBName != "" {
				cmd += fmt.Sprintf("--nsInclude %v.* ", overrideDBName)
			} else {
				cmd += fmt.Sprintf("--nsInclude %v.* ", target.Database)
			}
		}
	}

	if target.Params != "" {
		cmd += fmt.Sprintf("%v ", target.Params)
	}
	return cmd
}

func BuildUri(target config.Target) string {
	if target.Username != "" && target.Password != "" {
		return fmt.Sprintf("mongodb://%v:%v@%v:%v", target.Username, target.Password, target.Host, target.Port)
	}
	return fmt.Sprintf("mongodb://%v:%v", target.Host, target.Port)
}
