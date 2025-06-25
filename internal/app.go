package internal

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

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
var logLevel string
var logFormat string
var enableConcurrency bool
var concurrecyLimit int

// Add this function to app.go
func createExecutor() gcp.CommandExecutor {
	log := logger.New(appID, "createExecutor")

	if enableConcurrency {
		log.DebugWithExtra("Creating concurrent executor", map[string]any{
			"maxConcurrent": concurrecyLimit,
		})
		return gcp.NewConcurrentExecutor(concurrecyLimit)
	} else {
		log.Debug("Creating sequential executor")
		return &gcp.GCloudExecutor{} // Your original executor
	}
}

func Run(ctx context.Context) error {
	cli.Init(appID, shortDesc, longDesc)
	_ = cli.AddCommand("version", "Get the application version and Git commit SHA", logVersionDetails)
	_ = cli.AddCommand("check-health", "Check if we have the required tools installed", checkHealth)
	_ = cli.AddCommand("delete", "Delete all resources from a given folder", deleteResources)
	_ = cli.AddCommand("print", "Print the resource tree", printTree)
	cli.AssignStringFlag(&rootFolderId, "folder-id", "", "Root folder id to start from")
	cli.AssignStringFlag(&logLevel, "log-level", "info", "Log level (trace, debug, info, warn, error, fatal, panic)")
	cli.AssignStringFlag(&logFormat, "log-format", "pretty", "Log format (pretty, json)")
	cli.AssignBoolFlag(&dryRun, "dry-run", false, "Dry run mode")
	cli.AssignBoolFlag(&enableConcurrency, "concurrency", false, "Enable concurrency")
	cli.AssignIntFlag(&concurrecyLimit, "concurrency-limit", 5, "Concurrency limit")

	return cli.Run(ctx)
} // Updated helper function with format support

func initLogger(level string) error {
	if valid := validateLogLevel(level); !valid {
		return fmt.Errorf("invalid log level: %s", logLevel)
	}
	if valid := validateLogFormat(logFormat); !valid {
		return fmt.Errorf("invalid log format: %s", logFormat)
	}
	logger.Init(logger.Config{
		Level:  logLevel,
		Source: appID,
		Format: logFormat,
	})

	return nil
}

func validateLogLevel(level string) bool {
	validLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	level = strings.ToLower(level)
	return slices.Contains(validLevels, level)
}
func validateLogFormat(format string) bool {
	validFormats := []string{"pretty", "json"}
	format = strings.ToLower(format)
	return slices.Contains(validFormats, format)
}

func checkHealth(rootCtx context.Context) {
	_ = initLogger("info")
	executor := createExecutor()
	gcp.CheckHealth(rootCtx, executor)
}

func printTree(rootCtx context.Context) {
	_ = initLogger(logLevel)
	ctx, cancelFunc := context.WithCancel(rootCtx)
	defer cancelFunc()

	log := logger.New(appID, "printTree")

	if rootFolderId == "" {
		log.Error("rootFolderId is empty")
		return
	}

	executor := createExecutor()
	tree := getStructure(ctx, rootFolderId, executor)
	tree.Print()

}

func deleteResources(rootCtx context.Context) {
	_ = initLogger(logLevel)
	ctx, cancelFunc := context.WithCancel(rootCtx)
	defer cancelFunc()

	log := logger.New(appID, "deleteResources")

	if rootFolderId == "" {
		log.Error("rootFolderId is empty")
		return
	}

	executor := createExecutor()
	tree := getStructure(ctx, rootFolderId, executor)

	tree.Print()

	traversed := tree.PostOrderTraversal(tree.Root)
	log.DebugWithExtra("traversed", map[string]any{
		"traversed": traversed,
	})

	projects := make([]models.Entry, 0)
	folders := make([]models.Entry, 0)

	for _, entry := range traversed {
		switch entry.Type {
		case models.EntryTypeProject:
			projects = append(projects, entry)
		case models.EntryTypeFolder:
			folders = append(folders, entry)
		}
	}

	if enableConcurrency {
		var wg sync.WaitGroup

		for _, project := range projects {
			wg.Add(1)
			go func(p models.Entry) {
				defer wg.Done()
				err := gcp.DeleteProject(ctx, p.Id, dryRun, executor)
				if err != nil {
					log.Error("Failed to delete project", err)
				}
			}(project)
		}
		wg.Wait()
	} else {
		for _, project := range projects {
			err := gcp.DeleteProject(ctx, project.Id, dryRun, executor)
			if err != nil {
				log.Error("Failed to delete project", err)
			}
		}
	}

	for _, folder := range folders {
		err := gcp.DeleteFolder(ctx, folder.Id, dryRun, executor)
		if err != nil {
			log.Error("Failed to delete folder", err)
		}
	}

}

func logVersionDetails(_ context.Context) {
	_ = initLogger("info")
	log := logger.New(appID, "logVersionDetails")
	log.Info(fmt.Sprintf("AppVersion=%s, GitCommit=%s", version.AppVersion, version.GitCommit))
}
