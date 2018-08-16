package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Plan struct {
	Name      string    `yaml:"name"`
	Target    Target    `yaml:"target"`
	Scheduler Scheduler `yaml:"scheduler"`
	S3        *S3       `yaml:"s3"`
	GCloud    *GCloud   `yaml:"gcloud"`
	Azure     *Azure    `yaml:"azure"`
	SFTP      *SFTP     `yaml:"sftp"`
	SMTP      *SMTP     `yaml:"smtp"`
	Slack     *Slack    `yaml:"slack"`
}

type Target struct {
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Params   string `yaml:"params"`
}

type Scheduler struct {
	Cron      string `yaml:"cron"`
	Retention int    `yaml:"retention"`
	Timeout   int    `yaml:"timeout"`
}

type S3 struct {
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"accessKey"`
	API       string `yaml:"api"`
	SecretKey string `yaml:"secretKey"`
	URL       string `yaml:"url"`
}

type GCloud struct {
	Bucket      string `yaml:"bucket"`
	KeyFilePath string `yaml:"keyFilePath"`
}

type Azure struct {
	ContainerName string `yaml:"containerName"`
	ConnectionString string `yaml:"connectionString"`
}

type SFTP struct {
	Dir      string `yaml:"dir"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
}

type SMTP struct {
	Server   string   `yaml:"server"`
	Port     string   `yaml:"port"`
	Password string   `yaml:"password"`
	Username string   `yaml:"username"`
	From     string   `yaml:"from"`
	To       []string `yaml:"to"`
}

type Slack struct {
	URL      string `yaml:"url"`
	Channel  string `yaml:"channel"`
	Username string `yaml:"username"`
	WarnOnly bool   `yaml:"warnOnly"`
}

func LoadPlan(dir string, name string) (Plan, error) {
	plan := Plan{}
	planPath := ""
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, name+".yml") || strings.Contains(path, name+".yaml") {
			planPath = path
		}
		return nil
	})

	if err != nil {
		return plan, errors.Wrapf(err, "Reading from %v failed", dir)
	}

	if len(planPath) < 1 {
		return plan, errors.Errorf("Plan %v not found", name)
	}

	data, err := ioutil.ReadFile(planPath)
	if err != nil {
		return plan, errors.Wrapf(err, "Reading %v failed", planPath)
	}

	if err := yaml.Unmarshal(data, &plan); err != nil {
		return plan, errors.Wrapf(err, "Parsing %v failed", planPath)
	}
	_, filename := filepath.Split(planPath)
	plan.Name = strings.TrimSuffix(filename, filepath.Ext(filename))

	return plan, nil
}

func LoadPlans(dir string) ([]Plan, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "yml") || strings.Contains(path, "yaml") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "Reading from %v failed", dir)
	}

	plans := make([]Plan, 0)

	for _, path := range files {
		var plan Plan
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "Reading %v failed", path)
		}

		if err := yaml.Unmarshal(data, &plan); err != nil {
			return nil, errors.Wrapf(err, "Parsing %v failed", path)
		}
		_, filename := filepath.Split(path)
		plan.Name = strings.TrimSuffix(filename, filepath.Ext(filename))

		duplicate := false
		for _, p := range plans {
			if p.Name == plan.Name {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}

		plans = append(plans, plan)

	}
	if len(plans) < 1 {
		return nil, errors.Errorf("No backup plans found in %v", dir)
	}

	return plans, nil
}
