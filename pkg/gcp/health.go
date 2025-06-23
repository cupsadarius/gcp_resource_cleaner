package gcp

import (
	"context"
	"os/exec"
	"strings"

	"github.com/cupsadarius/gcp-resource-cleaner/pkg/logger"
)

func CheckHealth(rootCtx context.Context) {
	ctx, cancelFunc := context.WithCancel(rootCtx)
	defer cancelFunc()

	log := logger.New("gcp", "CheckHealth")
	cmd := exec.CommandContext(ctx, "gcloud", "version")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to run command", err)
	}

	if len(out) == 0 {
		log.Error("Gcloud command returned no output")
	}

	log.DebugWithExtra("Gcloud command output", map[string]any{
		"output": strings.Split(strings.Trim(string(out), "\n"), "\n"),
	})
}
