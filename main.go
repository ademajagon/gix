package main

import "github.com/ademajagon/gix/cmd"

func main() {
	v := cmd.Version()
	storedVersion = v
	go runCheckpoint(v)

	cmd.SetUpdateNotice(showCheckpointResult)
	cmd.Execute()
}
