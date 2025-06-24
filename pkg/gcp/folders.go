package gcp

import (
	"context"
	"strings"

	"github.com/cupsadarius/gcp_resource_cleaner/models"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
)

func GetFolders(rootCtx context.Context, rootFolderId string, executor CommandExecutor) ([]models.Entry, error) {
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
			"csv[no-heading](ID,DISPLAY_NAME)",
		},
	})
	out, err := executor.ExecuteCommand(ctx, "gcloud", "resource-manager", "folders", "list", "--folder", rootFolderId, "--format", "csv[no-heading](ID,DISPLAY_NAME)")
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
	result := make([]models.Entry, 0)
	trimmed := strings.Trim(string(out), "\n")
	for line := range strings.SplitSeq(trimmed, "\n") {
		if line != "" {
			vals := strings.Split(strings.Trim(line, "\n"), ",")
			entry := models.NewEntry(vals[0], vals[1], models.EntryTypeFolder)

			result = append(result, *entry)
		}
	}

	log.DebugWithExtra("Gcloud command output", map[string]any{
		"rootFolderId": rootFolderId,
		"output":       result,
	})

	return result, nil
}

func DeleteFolder(rootCtx context.Context, folderId string, dryRun bool, executor CommandExecutor) error {
	ctx, cancelFunc := context.WithCancel(rootCtx)
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

	out, err := executor.ExecuteCommand(ctx, "gcloud", "resource-manager", "folders", "delete", folderId, "--quiet")
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
