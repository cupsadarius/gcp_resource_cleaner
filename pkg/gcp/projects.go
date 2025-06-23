package gcp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
)

func GetProjects(rootCtx context.Context, rootFolderId string) ([]string, error) {
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
			"csv[no-heading](projectId)",
		},
	})
	cmd := exec.CommandContext(ctx, "gcloud", "projects", "list", "--filter", fmt.Sprintf("parent.id:%s", rootFolderId), "--format", "csv[no-heading](projectId)")
	out, err := cmd.CombinedOutput()
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

func DeleteProject(rootCtx context.Context, projectId string, dryRun bool) error {
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

	cmd := exec.CommandContext(ctx, "gcloud", "projects", "delete", projectId, "--quiet")
	out, err := cmd.CombinedOutput()
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
