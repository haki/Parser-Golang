package utilities

import (
	"Parser-Golang/services"
	"github.com/carlescere/scheduler"
)

func RegisterSchedulerFuncs() {
	scheduler.Every().Sunday().At("02:30").Run(services.AddNewComparisons)
	scheduler.Every().Day().At("00:01").Run(services.UpdateGitData)
}
