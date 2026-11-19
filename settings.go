package main

import "flag"

type Settings struct {
	Generate    bool
	LocalUpload bool
}

var GlobalSettings = &Settings{
	Generate:    false,
	LocalUpload: false,
}

func applyFlags() {
	// TODO
	// Flags/Settings global struct?
	// Add Dev mode flag
	generateFlag := flag.Bool("gen", false, "Generate typescript types to ./frontend")
	localUploadFlag := flag.Bool("localUpload", false, "Upload files to local dir instead of sftp server")
	flag.Parse()

	GlobalSettings.Generate = *generateFlag
	GlobalSettings.LocalUpload = *localUploadFlag
}
