package mongodump

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
	_ "github.com/codeskyblue/go-sh"
	"github.com/stefanprodan/mgob/config"
	_ "os/exec"
	"time"
	_ "time"
)

func Dump(plan config.Plan, conf *config.AppConfig) error {

	err := sh.Command("mkdir", fmt.Sprintf("%v/%v", conf.StoragePath, plan.Name)).Run()
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("mongodump --archive=%v/%v/syros.gz --gzip --host %v --port %v --db %v",
		conf.StoragePath, plan.Name, plan.Target.Host, plan.Target.Port, plan.Target.Database)

	out, err := sh.Command("/bin/sh", "-c", cmd).SetTimeout(time.Duration(plan.Scheduler.Timeout) * time.Minute).Output()

	fmt.Printf("output:(%s), err(%v)\n", string(out), err)

	if err != nil {
		return err
	}

	return nil
}
