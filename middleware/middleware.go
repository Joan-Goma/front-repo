package middleware

import (
	"fmt"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"neft.web/errorController"
)

func PrintStats() {
	ticker := time.NewTicker(2 * time.Minute)

	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			before, err := cpu.Get()
			if err != nil {
				errorController.ErrorLogger.Printf("%s\n", err)
				return
			}
			time.Sleep(time.Duration(1) * time.Second)
			after, err := cpu.Get()
			if err != nil {
				errorController.ErrorLogger.Printf("%s\n", err)
				return
			}
			memory, err := memory.Get()
			if err != nil {
				errorController.ErrorLogger.Printf("%s\n", err)
				return
			}

			total := float64(after.Total - before.Total)

			textPre := ("Printing stats\n-----os------stats---")

			currentCPU := fmt.Sprintf("\ncpu system:  %s %%", fmt.Sprintf("%.2f", float64(after.System-before.System)/total*100))

			totalCPU := fmt.Sprintf("\ncpu idle:    %s %%", fmt.Sprintf("%.2f", float64(after.Idle-before.Idle)/total*100))

			currentRAM := fmt.Sprintf("\nmemory used: %d  mb", memory.Used/1024/1024)

			freeRAM := fmt.Sprintf("\nmemory free: %d  mb", memory.Free/1024/1024)

			textPost := ("\n--------------------")
			errorController.DebugLogger.Println(textPre + currentCPU + totalCPU + currentRAM + freeRAM + textPost)

		case <-quit:
			ticker.Stop()
			return
		}
	}
}
