package backup

import (
	"fmt"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/config"
)

func s3Upload(file string, plan config.Plan) (string, error) {

	register := fmt.Sprintf("mc config host add %v %v %v %v --api %v",
		plan.Name, plan.S3.URL, plan.S3.AccessKey, plan.S3.SecretKey, plan.S3.API)

	result, err := sh.Command("/bin/sh", "-c", register).CombinedOutput()
	output := ""
	if len(result) > 0 {
		output = strings.Replace(string(result), "\n", " ", -1)
	}
	if err != nil {
		return "", errors.Wrapf(err, "mc config host for plan %v failed %s", plan.Name, output)
	}

	upload := fmt.Sprintf("mc --quiet cp %v %v/%v",
		file, plan.Name, plan.S3.Bucket)

	result, err = sh.Command("/bin/sh", "-c", upload).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	output = ""
	if len(result) > 0 {
		output = strings.Replace(string(result), "\n", " ", -1)
	}

	if err != nil {
		return "", errors.Wrapf(err, "S3 uploading %v to %v/%v failed %v", file, plan.Name, plan.S3.Bucket, output)
	}

	if strings.Contains(output, "<ERROR>") {
		return "", errors.Errorf("S3 upload failed %v", output)
	}

	return strings.Replace(output, "\n", " ", -1), nil
}
