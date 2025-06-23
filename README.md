# GCP Resource Cleaner

A robust CLI tool for recursively deleting GCP folders and projects in a safe, bottom-up approach. The tool builds a complete resource tree starting from a specified root folder, then systematically deletes resources from the leaves up to prevent dependency conflicts.

## How It Works

The GCP Resource Cleaner follows a **bottom-up deletion strategy** to safely remove GCP resources:

1. **Discovery Phase**: Recursively scans the GCP resource hierarchy starting from the provided folder ID
2. **Tree Building**: Constructs an in-memory tree representation of all folders and projects
3. **Post-Order Traversal**: Deletes resources from the deepest level first, working up to the root
4. **Safe Deletion**: Ensures child resources are removed before parent folders to avoid dependency errors

This approach prevents the common issue of trying to delete folders that still contain active projects or subfolders.

## Features

- **Recursive Resource Discovery**: Automatically finds all projects and subfolders within a given folder hierarchy
- **Safe Bottom-Up Deletion**: Deletes resources in the correct order to avoid dependency conflicts
- **Dry-Run Mode**: Preview what would be deleted without making any actual changes
- **Comprehensive Logging**: Detailed debug output for troubleshooting and verification
- **Health Checks**: Validates that required tools (gcloud CLI) are properly configured
- **Signal Handling**: Graceful shutdown on interruption

## Prerequisites

Before using this tool, ensure you have:

- **gcloud CLI** installed and authenticated
- **Appropriate GCP Permissions**:
  - `resourcemanager.folders.delete` on target folders
  - `resourcemanager.projects.delete` on target projects
  - `resourcemanager.folders.list` and `resourcemanager.projects.list` for discovery
- **A GCP Folder ID** as the starting point for deletion

## Installation

### Install from source:
```bash
go install github.com/cupsadarius/gcp_resource_cleaner@latest
```

### Build locally:
```bash
git clone https://github.com/cupsadarius/gcp_resource_cleaner.git
cd gcp_resource_cleaner
make mod-tidy
make build
```

## Usage

### Check System Health
Verify that gcloud CLI is properly installed and configured:
```bash
gcp_resource_cleaner check-health
```

### Preview Deletion (Dry Run)
**Always start with a dry run** to see what would be deleted:
```bash
gcp_resource_cleaner delete --folder-id <folder-id> --dry-run
```

Example:
```bash
gcp_resource_cleaner delete --folder-id 123456789012 --dry-run
```

### Execute Deletion
After reviewing the dry-run output, proceed with actual deletion:
```bash
gcp_resource_cleaner delete --folder-id <folder-id>
```

### Get Version Information
```bash
gcp_resource_cleaner version
```

## Command Reference

| Command | Description | Flags |
|---------|-------------|-------|
| `check-health` | Validates gcloud CLI installation and authentication | None |
| `delete` | Recursively deletes folders and projects | `--folder-id` (required), `--dry-run` (optional) |
| `version` | Shows application version and Git commit SHA | None |

## Example Workflow

```bash
# 1. Verify system health
gcp_resource_cleaner check-health

# 2. Preview what would be deleted
gcp_resource_cleaner delete --folder-id 123456789012 --dry-run

# 3. Review the output carefully, then execute
gcp_resource_cleaner delete --folder-id 123456789012
```

## Safety Features

- **Dry-run mode** prevents accidental deletions
- **Bottom-up traversal** ensures safe deletion order
- **Comprehensive logging** provides full visibility into operations
- **Signal handling** allows graceful cancellation
- **Error handling** stops on failures to prevent partial deletions

## Architecture

The tool is structured with clear separation of concerns:

- `internal/`: Core application logic and orchestration
- `models/`: Data structures for tree representation and resource entries
- `pkg/cli/`: Command-line interface handling
- `pkg/gcp/`: GCP API interactions via gcloud CLI
- `pkg/logger/`: Structured logging with zerolog

## Troubleshooting

### Common Issues

**"Failed to get projects/folders"**
- Verify gcloud authentication: `gcloud auth list`
- Check folder ID exists: `gcloud resource-manager folders describe <folder-id>`
- Ensure proper IAM permissions

**"Permission denied" errors**
- Verify you have `resourcemanager.folders.delete` and `resourcemanager.projects.delete` permissions
- Check if folders/projects have any protection policies

**"Folder not empty" errors**
- This shouldn't happen with proper post-order traversal
- Check logs for failed deletions of child resources
- Some resources may require manual cleanup (e.g., billing accounts, liens)

### Debug Logging

The tool provides detailed debug output. Enable verbose logging to troubleshoot issues:
```bash
# Logs are automatically set to debug level and show:
# - Resource discovery operations
# - Tree structure building
# - Deletion sequence and results
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Important Notes

⚠️ **Resource Deletion Behavior:**
- **Projects**: Enter a 30-day pending deletion period and can be recovered via the [Cloud Resource Manager](https://console.cloud.google.com/cloud-resource-manager?pendingDeletion=true)
- **Folders**: Are permanently deleted and **cannot be recovered**
- Always use `--dry-run` first and verify the output before proceeding with actual deletion
