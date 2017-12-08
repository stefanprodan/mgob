package backup

import (
	"fmt"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/config"
)

func gCloudUpload(file string, plan config.Plan) (string, error) {

	register := fmt.Sprintf("gcloud auth activate-service-account --key-file=%v",
		plan.GCloud.KeyFilePath)

	_, err := sh.Command("/bin/sh", "-c", register).CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "gcloud auth for plan %v failed", plan.Name)
	}

	upload := fmt.Sprintf("gsutil cp %v gs://%v",
		file, plan.GCloud.Bucket)

	result, err := sh.Command("/bin/sh", "-c", upload).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	output := ""
	if len(result) > 0 {
		output = strings.Replace(string(result), "\n", " ", -1)
	}

	if err != nil {
		return "", errors.Wrapf(err, "GCloud uploading %v to gs://%v failed %v", file, plan.GCloud.Bucket, output)
	}

	if strings.Contains(output, "<ERROR>") {
		return "", errors.Errorf("GCloud upload failed %v", output)
	}

	return strings.Replace(output, "\n", " ", -1), nil
}
