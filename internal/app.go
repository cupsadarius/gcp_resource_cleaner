package internal

import (
	"context"
	"fmt"

	"github.com/cupsadarius/gcp_resource_cleaner/models"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/cli"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/gcp"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/version"
)

const appID = "gcp_resource_cleaner"
const shortDesc = "CLI Tool for cleaning up GCP resources"
const longDesc = `gcp_resource_cleaner is a CLI tool for cleaning up GCP resources.

Since GCP prevents you from deleting resources that are in use,
this tool recursively traverses the GCP resource tree and deletes
all resources from bottom up given a starting folder id.`

var rootFolderId string
var dryRun bool

func Run(ctx context.Context) error {
	cli.Init(appID, shortDesc, longDesc)
	_ = cli.AddCommand("version", "Get the application version and Git commit SHA", logVersionDetails)
	_ = cli.AddCommand("check-health", "Check if we have the required tools installed", gcp.CheckHealth)
	_ = cli.AddCommand("delete", "Delete all resources from a given folder", deleteResources)
	cli.AssignStringFlag(&rootFolderId, "folder-id", "", "Root folder id to start from")
	cli.AssignBoolFlag(&dryRun, "dry-run", false, "Dry run mode")

	logger.Init(logger.Config{
		Level:  "info",
		Source: appID,
		Format: "pretty",
	})
	return cli.Run(ctx)
}

func deleteResources(rootCtx context.Context) {
	ctx, cancelFunc := context.WithCancel(rootCtx)
	defer cancelFunc()

	log := logger.New(appID, "deleteResources")

	if rootFolderId == "" {
		log.Error("rootFolderId is empty")
		return
	}

	tree := getStructure(ctx, rootFolderId)
	traversed := tree.PostOrderTraversal(tree.Root)
	log.DebugWithExtra("traversed", map[string]any{
		"traversed": traversed,
	})

	for _, entry := range traversed {
		switch entry.Type {
		case models.EntryTypeProject:
			err := gcp.DeleteProject(ctx, entry.Id, dryRun)
			if err != nil {
				log.Error("Failed to delete project", err)
			}
		case models.EntryTypeFolder:
			err := gcp.DeleteFolder(ctx, entry.Id, dryRun)
			if err != nil {
				log.Error("Failed to delete folder", err)
			}
		}
	}

}

func logVersionDetails(_ context.Context) {
	log := logger.New(appID, "logVersionDetails")
	log.Info(fmt.Sprintf("AppVersion=%s, GitCommit=%s", version.AppVersion, version.GitCommit))
}
