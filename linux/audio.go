package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func playAudio() {
	audioPath := filepath.Join(os.TempDir(), "audio.mp3")
	os.WriteFile(audioPath, AudioBytes, 0666)

	mpg := GetMPGCommand()

	ctx, cancel := context.WithCancel(context.Background())
	linuxCancelFn = cancel
	out, err := exec.CommandContext(ctx, mpg, audioPath).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}

func GetMPGCommand() string {
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "mpg321"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "usr", "bin", "mpg321")
	}

	return cmdPath
}
