package gcp

import (
	"context"
	"strings"

	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
)

func CheckHealth(rootCtx context.Context, executor CommandExecutor) {
	ctx, cancelFunc := context.WithCancel(rootCtx)
	defer cancelFunc()

	log := logger.New("gcp", "CheckHealth")
	out, err := executor.ExecuteCommand(ctx, "gcloud", "version")

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
