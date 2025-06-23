package gcp

import (
	"context"
	"os/exec"
	"strings"

	"github.com/cupsadarius/gcp-resource-cleaner/pkg/logger"
)

func GetFolders(rootCtx context.Context, rootFolderId string) ([]string, error) {
	ctx, cancelFunc := context.WithCancel(rootCtx)
	defer cancelFunc()
	log := logger.New("gcp", "GetFolders")
	log.DebugWithExtra("getFolders", map[string]any{
		"cmd": "gcloud",
		"args": []string{
			"resource-manager",
			"folders",
			"list",
			"--folder",
			rootFolderId,
			"--format",
			"csv[no-heading](ID)",
		},
	})
	cmd := exec.CommandContext(ctx, "gcloud", "resource-manager", "folders", "list", "--folder", rootFolderId, "--format", "csv[no-heading](ID)")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to run command", err)
		return nil, err
	}

	if len(out) == 0 {
		log.DebugWithExtra("There are no folders in the given folder", map[string]any{
			"rootFolderId": rootFolderId,
		})
		return nil, nil
	}
	result := make([]string, 0)
	trimmed := strings.Trim(string(out), "\n")
	for line := range strings.SplitSeq(trimmed, "\n") {
		if line != "" {
			result = append(result, strings.Trim(line, "\n"))
		}
	}

	log.DebugWithExtra("Gcloud command output", map[string]any{
		"rootFolderId": rootFolderId,
		"output":       result,
	})

	return result, nil
}

func DeleteFolder(ctx context.Context, folderId string, dryRun bool) error {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	log := logger.New("gcp", "DeleteFolder")
	log.DebugWithExtra("DeleteFolder", map[string]any{
		"cmd": "gcloud",
		"args": []string{
			"resource-manager",
			"folders",
			"delete",
			folderId,
			"--quiet",
		},
	})

	if dryRun {
		return nil
	}

	cmd := exec.CommandContext(ctx, "gcloud", "resource-manager", "folders", "delete", folderId, "--quiet")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to run command", err)
		return err
	}

	if len(out) == 0 {
		log.DebugWithExtra("Gcloud command returned no output", map[string]any{
			"folderId": folderId,
		})
		return nil
	}
	log.DebugWithExtra("Gcloud command output", map[string]any{
		"folderId": folderId,
		"output":   strings.Split(strings.Trim(string(out), "\n"), "\n"),
	})

	return nil
}
