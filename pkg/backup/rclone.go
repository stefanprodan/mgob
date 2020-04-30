package backup

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"

	"github.com/stefanprodan/mgob/pkg/config"
)

func rcloneUpload(file string, plan config.Plan) (string, error) {

	fileName := filepath.Base(file)

	upload := fmt.Sprintf("rclone --config=\"%v\" copy %v %v:%v/%v",
		plan.Rclone.ConfigFilePath, plan.Name, plan.Rclone.Bucket, fileName)

	result, err := sh.Command("/bin/sh", "-c", upload).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	output := ""
	if len(result) > 0 {
		output = strings.Replace(string(result), "\n", " ", -1)
	}

	if err != nil {
		return "", errors.Wrapf(err, "Rclone uploading %v to %v:%v failed %v", file, plan.Name, plan.Rclone.Bucket, output)
	}

	return strings.Replace(output, "\n", " ", -1), nil
}
