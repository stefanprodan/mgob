package backup

import (
	"github.com/pkg/errors"
	"github.com/codeskyblue/go-sh"
	"strings"
)

func CheckMongodump() (string, error) {
	output, err := sh.Command("/bin/sh", "-c", "mongodump --version").CombinedOutput()
	if err != nil {
		ex := ""
		if len(output) > 0 {
			ex = strings.Replace(string(output), "\n", " ", -1)
		}
		return "", errors.Wrapf(err, "mongodump failed %v", ex)
	}

	return strings.Replace(string(output), "\n", " ", -1), nil
}

func CheckMinioClient() (string, error) {
	output, err := sh.Command("/bin/sh", "-c", "mc --version").CombinedOutput()
	if err != nil {
		ex := ""
		if len(output) > 0 {
			ex = strings.Replace(string(output), "\n", " ", -1)
		}
		return "", errors.Wrapf(err, "mc failed %v", ex)
	}

	return strings.Replace(string(output), "\n", " ", -1), nil
}
