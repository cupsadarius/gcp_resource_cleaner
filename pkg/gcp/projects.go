package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/cupsadarius/gcp_resource_cleaner/models"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
)

func GetProjects(rootCtx context.Context, rootFolderId string, executor CommandExecutor) ([]models.Entry, error) {
	ctx, cancelFunc := context.WithCancel(rootCtx)
	defer cancelFunc()

	log := logger.New("gcp", "GetProjects")
	log.DebugWithExtra("getProjects", map[string]any{
		"cmd": "gcloud",
		"args": []string{
			"projects",
			"list",
			"--filter",
			fmt.Sprintf("parent.id:%s", rootFolderId),
			"--format",
			"csv[no-heading](projectId,name)",
		},
	})
	out, err := executor.ExecuteCommand(ctx, "gcloud", "projects", "list", "--filter", fmt.Sprintf("parent.id:%s", rootFolderId), "--format", "csv[no-heading](projectId,name)")
	if err != nil {
		log.Error("Failed to run command", err)
		return nil, err
	}

	if len(out) == 0 {
		log.DebugWithExtra("There are no projects in the given folder", map[string]any{
			"rootFolderId": rootFolderId,
		})
		return nil, nil
	}
	result := make([]models.Entry, 0)
	trimmed := strings.Trim(string(out), "\n")
	for line := range strings.SplitSeq(trimmed, "\n") {
		if line != "" {
			vals := strings.Split(strings.Trim(line, "\n"), ",")
			entry := models.NewEntry(vals[1], vals[0], models.EntryTypeProject)
			result = append(result, *entry)
		}
	}

	log.DebugWithExtra("Gcloud command output", map[string]any{
		"rootFolderId": rootFolderId,
		"output":       result,
	})

	return result, nil
}

func DeleteProject(rootCtx context.Context, projectId string, dryRun bool, executor CommandExecutor) error {
	ctx, cancelFunc := context.WithCancel(rootCtx)
	defer cancelFunc()

	log := logger.New("gcp", "DeleteProject")
	log.DebugWithExtra("DeleteProject", map[string]any{
		"cmd": "gcloud",
		"args": []string{
			"projects",
			"delete",
			projectId,
			"--quiet",
		},
	})
	if dryRun {
		return nil
	}

	out, err := executor.ExecuteCommand(ctx, "gcloud", "projects", "delete", projectId, "--quiet")
	if err != nil {
		log.Error("Failed to run command", err)
		return err
	}

	if len(out) == 0 {
		log.DebugWithExtra("Gcloud command returned no output", map[string]any{
			"projectId": projectId,
		})
		return nil
	}
	log.DebugWithExtra("Gcloud command output", map[string]any{
		"projectId": projectId,
		"output":    strings.Split(strings.Trim(string(out), "\n"), "\n"),
	})

	return nil
}
