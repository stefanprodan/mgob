package backup

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"net/url"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"

	"github.com/stefanprodan/mgob/pkg/config"
)

func s3Upload(file string, plan config.Plan, useAwsCli bool) (string, error) {

	s3Url, err := url.Parse(plan.S3.URL)

	if err != nil {
		return "", errors.Wrapf(err, "invalid S3 url for plan %v: %s", plan.Name, plan.S3.URL)
	}

	if useAwsCli && strings.HasSuffix(s3Url.Hostname(), "amazonaws.com") {
		return awsUpload(file, plan)
	}

	return minioUpload(file, plan)
}

func awsUpload(file string, plan config.Plan) (string, error) {

	output := ""
	if len(plan.S3.AccessKey) > 0 && len(plan.S3.SecretKey) > 0 {
		// Let's use credentials given
		configure := fmt.Sprintf("aws configure set aws_access_key_id %v && aws configure set aws_secret_access_key %v",
			plan.S3.AccessKey, plan.S3.SecretKey)

		result, err := sh.Command("/bin/sh", "-c", configure).CombinedOutput()
		if len(result) > 0 {
			output += strings.Replace(string(result), "\n", " ", -1)
		}
		if err != nil {
			return "", errors.Wrapf(err, "aws configure for plan %v failed %s", plan.Name, output)
		}
	}

	fileName := filepath.Base(file)

	encrypt := ""
	if len(plan.S3.KmsKeyId) > 0 {
		encrypt = fmt.Sprintf(" --sse aws:kms --sse-kms-key-id %v", plan.S3.KmsKeyId)
	}

	storage := ""
	if len(plan.S3.StorageClass) > 0 {
		storage = fmt.Sprintf(" --storage-class %v", plan.S3.StorageClass)
	}

	upload := fmt.Sprintf("aws --quiet s3 cp %v s3://%v/%v%v%v",
		file, plan.S3.Bucket, fileName, encrypt, storage)

	result, err := sh.Command("/bin/sh", "-c", upload).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).CombinedOutput()
	if len(result) > 0 {
		output += strings.Replace(string(result), "\n", " ", -1)
	}
	if err != nil {
		return "", errors.Wrapf(err, "S3 uploading %v to %v/%v failed %v", file, plan.Name, plan.S3.Bucket, output)
	}

	if strings.Contains(output, "<ERROR>") {
		return "", errors.Errorf("S3 upload failed %v", output)
	}

	return strings.Replace(output, "\n", " ", -1), nil
}

func minioUpload(file string, plan config.Plan) (string, error) {

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

	fileName := filepath.Base(file)

	upload := fmt.Sprintf("mc --quiet cp %v %v/%v/%v",
		file, plan.Name, plan.S3.Bucket, fileName)

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
