package backup

import (
	"fmt"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"github.com/stefanprodan/mgob/config"
)

func azureUpload(file string, plan config.Plan) (string, error) {
	upload := fmt.Sprintf("az storage blob upload -c '%v' --file '%v' --name '%v' --connection-string '%v'",
		plan.Azure.ContainerName, file, file, plan.Azure.ConnectionString)

	result, err := sh.Command("/bin/sh", "-c", upload).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	output := ""
	if len(result) > 0 {
		output = strings.Replace(string(result), "\n", " ", -1)
	}

	if err != nil {
		return "", errors.Wrapf(err, "Azure uploading %v to %v failed %v", file, plan.Azure.ContainerName, output)
	}

	if strings.Contains(output, "<Error>") {
		return "", errors.Errorf("Azure upload failed %v", output)
	}

	return strings.Replace(output, "\n", " ", -1), nil
}
