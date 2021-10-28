package main

import (
	"log"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

func main() {

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}

	for _, s := range stat.CPUStats {

		fmt.Println(s.User)
		// s.User
		// s.Nice
		// s.System
		// s.Idle
		// s.IOWait
	}

	// stat.CPUStatAll
	// stat.CPUStats
	// stat.Processes
	// stat.BootTime

}

