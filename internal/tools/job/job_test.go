package job

import (
	"fmt"
	"github.com/reugn/go-quartz/quartz"
	"log"
	"testing"
)

func TestCron(t *testing.T) {
	cron := "0 30 9 * * ?"
	trigger, err := quartz.NewCronTrigger(cron)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(trigger.NextFireTime(1))
}
