package config

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Plan struct {
	Name      string    `yaml:"name"`
	Target    Target    `yaml:"target"`
	Scheduler Scheduler `yaml:"scheduler"`
}

type Target struct {
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
}

type Scheduler struct {
	Cron      string `yaml:"cron"`
	Retention int    `yaml:"retention"`
}

func LoadPlans(dir string) ([]Plan, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "yml") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "reading from %v failed", dir)
	}

	plans := make([]Plan, 0)

	for _, path := range files {
		var plan Plan
		if strings.Contains(path, "yml") {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, errors.Wrapf(err, "reading %v failed", path)
			}

			if err := yaml.Unmarshal(data, &plan); err != nil {
				return nil, errors.Wrapf(err, "parsering %v failed", path)
			}
			_, filename := filepath.Split(path)
			plan.Name = strings.TrimSuffix(filename, filepath.Ext(filename))
			plans = append(plans, plan)

			logrus.Infof("Plan %v loaded", filename)
		}
	}
	if len(plans) < 1 {
		return nil, errors.New(fmt.Sprintf("No backup plans found in %v", dir))
	}

	return plans, nil
}
