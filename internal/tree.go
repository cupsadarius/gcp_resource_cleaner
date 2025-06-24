package internal

import (
	"context"

	"github.com/cupsadarius/gcp_resource_cleaner/models"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/gcp"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
)

func getStructure(ctx context.Context, rootFolderId string, executor gcp.CommandExecutor) *models.Tree {
	tree := models.NewTree()
	rootEntry := models.NewEntry(rootFolderId, rootFolderId, models.EntryTypeFolder)

	tree.Root = getTree(ctx, *rootEntry, executor)

	return tree
}

func getTree(ctx context.Context, root models.Entry, executor gcp.CommandExecutor) *models.Node {
	log := logger.New(appID, "getStructure")
	log.DebugWithExtra("getStructure", map[string]any{
		"rootFolderId": rootFolderId,
	})
	projects, err := gcp.GetProjects(ctx, root.Id, executor)
	if err != nil {
		log.Error("Failed to get projects", err)
		return nil
	}

	node := models.NewNode(&root, projects)

	folders, err := gcp.GetFolders(ctx, root.Id, executor)
	if err != nil {
		log.Error("Failed to get folders", err)
		return node
	}
	for _, folder := range folders {
		node.Children = append(node.Children, getTree(ctx, folder, executor))
	}

	return node
}
