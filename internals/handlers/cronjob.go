package handlers

import (
	"IFEST/internals/services"
	"github.com/robfig/cron/v3"
	"log"
)

type CronJob struct {
	userDocService services.IUserDocService
	cronScheduler  *cron.Cron
}

func NewCronJob(userDocService services.IUserDocService) CronJob {
	return CronJob{
		userDocService: userDocService,
		cronScheduler:  cron.New(),
	}
}

func (cj *CronJob) Start() {
	_, err := cj.cronScheduler.AddFunc("*/10 * * * *", func() {
		log.Println("CronJob: delete expired access...")
		err := cj.userDocService.DeleteExpired()
		if err != nil {
			log.Printf("CronJob Error: %v\n", err)
		} else {
			log.Println("CronJob: deleted successfully")
		}
	})
	if err != nil {
		log.Fatalf("failed to add cron: %v", err)
	}

	cj.cronScheduler.Start()
	log.Println("CronJob: Scheduler started.")
}

func (cj *CronJob) Stop() {
	ctx := cj.cronScheduler.Stop()
	<-ctx.Done()
	log.Println("CronJob: Scheduler stopped.")
}
