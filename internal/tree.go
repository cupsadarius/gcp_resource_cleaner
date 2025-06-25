package internal

import (
	"context"
	"sync"

	"github.com/cupsadarius/gcp_resource_cleaner/models"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/gcp"
	"github.com/cupsadarius/gcp_resource_cleaner/pkg/logger"
)

func getStructure(ctx context.Context, rootFolderId string, executor gcp.CommandExecutor) *models.Tree {
	tree := models.NewTree()
	rootEntry := models.NewEntry(rootFolderId, rootFolderId, models.EntryTypeFolder)

	if enableConcurrency {
		tree.Root = getTreeWithConcurrentSubfolders(ctx, *rootEntry, executor)
	} else {
		// EXISTING: Use your original sequential version
		tree.Root = getTree(ctx, *rootEntry, executor)
	}

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

func getTreeWithConcurrentSubfolders(ctx context.Context, root models.Entry, executor gcp.CommandExecutor) *models.Node {
	log := logger.New(appID, "getTreeWithConcurrentSubfolders")

	// Get projects and folders for current folder (sequential)
	projects, err := gcp.GetProjects(ctx, root.Id, executor)
	if err != nil {
		log.Error("Failed to get projects", err)
		return nil
	}

	folders, err := gcp.GetFolders(ctx, root.Id, executor)
	if err != nil {
		log.Error("Failed to get folders", err)
		return nil
	}

	// Create node
	node := models.NewNode(&root, projects)

	// Process subfolders concurrently (but still recursively sequential)
	if len(folders) > 0 {
		var wg sync.WaitGroup
		children := make([]*models.Node, len(folders))

		// Create one goroutine per subfolder
		for i, folder := range folders {
			wg.Add(1)

			go func(index int, folderEntry models.Entry) {
				defer wg.Done()

				log.DebugWithExtra("Processing subfolder", map[string]any{
					"index":    index,
					"folderId": folderEntry.Id,
				})

				// Recursive call (still sequential within each subtree)
				children[index] = getTreeWithConcurrentSubfolders(ctx, folderEntry, executor)
			}(i, folder)
		}

		wg.Wait()

		// Add non-nil children
		for _, child := range children {
			if child != nil {
				node.Children = append(node.Children, child)
			}
		}
	}

	return node
}
