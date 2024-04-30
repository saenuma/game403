package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/saenuma/Game403/g403l"
)

func playAudio() {
	audioPath := filepath.Join(os.TempDir(), "audio.mp3")
	os.WriteFile(audioPath, g403l.AudioBytes, 0666)

	mpg := GetMPGCommand()

	ctx, cancel := context.WithCancel(context.Background())
	linuxCancelFn = cancel
	out, err := exec.CommandContext(ctx, mpg, audioPath).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
	}
}

func GetMPGCommand() string {
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = "madplay"
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "usr", "bin", "madplay")
	}

	return cmdPath
}
