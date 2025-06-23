package internal

import (
	"context"

	"github.com/cupsadarius/gcp_resource_cleaner/models"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/gcp"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
)

func getStructure(ctx context.Context, rootFolderId string, executor gcp.CommandExecutor) *models.Tree {
	tree := models.NewTree()

	tree.Root = getTree(ctx, rootFolderId, executor)

	return tree
}

func getTree(ctx context.Context, rootFolderId string, executor gcp.CommandExecutor) *models.Node {
	log := logger.New(appID, "getStructure")
	log.DebugWithExtra("getStructure", map[string]any{
		"rootFolderId": rootFolderId,
	})
	projects, err := gcp.GetProjects(ctx, rootFolderId, executor)
	if err != nil {
		log.Error("Failed to get projects", err)
		return nil
	}

	node := models.NewNode(rootFolderId, projects)

	folders, err := gcp.GetFolders(ctx, rootFolderId, executor)
	if err != nil {
		log.Error("Failed to get folders", err)
		return node
	}
	for _, folder := range folders {
		node.Children = append(node.Children, getTree(ctx, folder, executor))
	}

	return node
}
