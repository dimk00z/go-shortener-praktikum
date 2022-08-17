package util

import "fmt"

const notGiven = "N/A"

func PrintBulidInfo(
	buildVersion string,
	buildDate string,
	buildCommit string) {

	if buildVersion == "" {
		buildVersion = notGiven
	}
	if buildDate == "" {
		buildDate = notGiven
	}
	if buildCommit == "" {
		buildCommit = notGiven
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
