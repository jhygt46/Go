package main

import (
	"log"
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

func main() {

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}

	for _, s := range stat.CPUStats {

		fmt.Println("User", s.User)
		fmt.Println("Nice", s.Nice)
		fmt.Println("System", s.System)
		fmt.Println("Idle", s.Idle)
		fmt.Println("IOWait", s.IOWait)

	}

	// stat.CPUStatAll
	// stat.CPUStats
	// stat.Processes
	// stat.BootTime

}

