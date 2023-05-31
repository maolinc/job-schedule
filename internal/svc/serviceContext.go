package svc

import (
	"job/internal/config"
	"job/internal/tools/job"
	"job/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config config.Config

	JobRecordModel   model.JobRecordModel
	JobScheduleModel model.JobScheduleModel
	Scheduler        *job.Scheduler
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, _ := gorm.Open(mysql.Open(c.DB.DataSource), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	recordModel := model.NewJobRecordModel(db)
	scheduleModel := model.NewJobScheduleModel(db)

	sc := &ServiceContext{
		Config: c,

		JobRecordModel:   recordModel,
		JobScheduleModel: scheduleModel,

		Scheduler: job.NewScheduler(&job.DB{
			JobRecord:   recordModel,
			JobSchedule: scheduleModel,
		}),
	}

	return sc
}
