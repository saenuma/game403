package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

func playAudio() {
	audioPath := filepath.Join(os.TempDir(), "audio.mp3")
	os.WriteFile(audioPath, AudioBytes, 0666)

	ctx, cancel := context.WithCancel(context.Background())
	linuxCancelFn = cancel
	cmd := exec.CommandContext(ctx, "mpg321", audioPath)
	cmd.Run()
}
