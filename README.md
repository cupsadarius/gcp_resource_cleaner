# GCP Resource Cleaner

A robust CLI tool for recursively deleting GCP folders and projects in a safe, bottom-up approach. The tool builds a complete resource tree starting from a specified root folder, then systematically deletes resources from the leaves up to prevent dependency conflicts.

## How It Works

The GCP Resource Cleaner follows a **bottom-up deletion strategy** to safely remove GCP resources:

1. **Discovery Phase**: Recursively scans the GCP resource hierarchy starting from the provided folder ID
2. **Tree Building**: Constructs an in-memory tree representation of all folders and projects
3. **Visualization**: Displays the discovered resource hierarchy in an easy-to-read tree format
4. **Post-Order Traversal**: Deletes resources from the deepest level first, working up to the root
5. **Safe Deletion**: Ensures child resources are removed before parent folders to avoid dependency errors

This approach prevents the common issue of trying to delete folders that still contain active projects or subfolders.

## Features

- **Recursive Resource Discovery**: Automatically finds all projects and subfolders within a given folder hierarchy
- **Pretty Tree Visualization**: Visual representation of the resource hierarchy before deletion
- **Safe Bottom-Up Deletion**: Deletes resources in the correct order to avoid dependency conflicts
- **Concurrent Processing**: Optional concurrent discovery and deletion for improved performance on large hierarchies
- **Dry-Run Mode**: Preview what would be deleted without making any actual changes
- **Configurable Logging**: Adjustable log levels from silent operation to detailed debugging
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

### Print Resource Tree
Discover and visualize the resource hierarchy without any deletion operations:
```bash
gcp_resource_cleaner print --folder-id <folder-id>
```

### Preview Deletion Plan
View the resource hierarchy and deletion plan without making any changes:
```bash
gcp_resource_cleaner delete --folder-id <folder-id> --dry-run
```

### Preview Deletion with Detailed Logging
**Always start with a dry run** to see what would be deleted with verbose output:
```bash
gcp_resource_cleaner delete --folder-id <folder-id> --dry-run --log-level debug
```

Example:
```bash
gcp_resource_cleaner delete --folder-id <foleder-id> --dry-run --log-level debug
```

### Execute Deletion with Minimal Logging
After reviewing the dry-run output, proceed with actual deletion:
```bash
gcp_resource_cleaner delete --folder-id <folder-id> --log-level warn
```

### Concurrent Processing
For large folder hierarchies, enable concurrent processing to improve performance:

```bash
# Enable concurrency with default limit (5 concurrent operations)
gcp_resource_cleaner delete --folder-id <folder-id> --concurrency --dry-run

# Enable concurrency with custom limit
gcp_resource_cleaner delete --folder-id <folder-id> --concurrency --concurrency-limit 10 --dry-run

# High concurrency for very large hierarchies
gcp_resource_cleaner delete --folder-id <folder-id> --concurrency --concurrency-limit 20 --dry-run
```

### Advanced Usage Examples
```bash
# Just view the resource tree structure
gcp_resource_cleaner print --folder-id <folder-id> --log-level info

# View tree with detailed discovery logging and concurrency
gcp_resource_cleaner print --folder-id <folder-id> --log-level debug --concurrency --concurrency-limit 8

# Detailed debugging with pretty console output and concurrency
gcp_resource_cleaner delete --folder-id <folder-id> --dry-run --log-level trace --log-format pretty --concurrency

# Production deletion with JSON logging and moderate concurrency for log aggregation
gcp_resource_cleaner delete --folder-id <folder-id> --log-level error --log-format json --concurrency --concurrency-limit 10

# Silent operation with maximum safe concurrency (errors only)
gcp_resource_cleaner delete --folder-id <folder-id> --log-level error --concurrency --concurrency-limit 15

# Maximum verbosity with concurrency for troubleshooting large hierarchies
gcp_resource_cleaner delete --folder-id <folder-id> --dry-run --log-level trace --concurrency --concurrency-limit 5
```

### Get Version Information
```bash
gcp_resource_cleaner version
```

## Command Reference

| Command | Description | Flags |
|---------|-------------|-------|
| `check-health` | Validates gcloud CLI installation and authentication | `--log-level`, `--log-format` |
| `print` | Displays the resource tree structure without any deletion operations | `--folder-id` (required), `--log-level`, `--log-format`, `--concurrency`, `--concurrency-limit` |
| `delete` | Recursively deletes folders and projects | `--folder-id` (required), `--dry-run`, `--log-level`, `--log-format`, `--concurrency`, `--concurrency-limit` |
| `version` | Shows application version and Git commit SHA | `--log-level`, `--log-format` |

## Flag Reference

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--folder-id` | string | "" | Root folder ID to start from (required for print and delete commands) |
| `--dry-run` | bool | false | Preview mode - shows what would be deleted without making changes (delete command only) |
| `--log-level` | string | "info" | Log verbosity level: trace, debug, info, warn, error, fatal, panic |
| `--log-format` | string | "pretty" | Log output format: pretty (human-readable) or json (machine-readable) |
| `--concurrency` | bool | false | Enable concurrent processing for improved performance |
| `--concurrency-limit` | int | 5 | Maximum number of concurrent operations (only applies when `--concurrency` is enabled) |


## Performance Optimization

### Sequential vs Concurrent Processing

**Sequential Processing (default)**:
- Processes one folder at a time
- Lower memory usage
- Easier to debug
- Recommended for small to medium hierarchies

**Concurrent Processing**:
- Processes multiple folders simultaneously
- Higher throughput for large hierarchies
- Increased memory usage
- May hit API rate limits if set too high

## Safety Features

- **Dry-run mode** prevents accidental deletions
- **Tree visualization** shows complete resource hierarchy before deletion
- **Bottom-up traversal** ensures safe deletion order
- **Configurable logging** provides appropriate verbosity for different use cases
- **Concurrent processing** with rate limiting to respect API limits
- **Signal handling** allows graceful cancellation
- **Error handling** stops on failures to prevent partial deletions

## Architecture

The tool is structured with clear separation of concerns:

- `internal/`: Core application logic and orchestration
- `models/`: Data structures for tree representation and resource entries
- `pkg/cli/`: Command-line interface handling with configurable logging
- `pkg/gcp/`: GCP API interactions via gcloud CLI with concurrent execution support
- `pkg/logger/`: Structured logging with zerolog (configurable levels and formats)

## Troubleshooting

### Common Issues

**"Failed to get projects/folders"**
- Verify gcloud authentication: `gcloud auth list`
- Check folder ID exists: `gcloud resource-manager folders describe <folder-id>`
- Ensure proper IAM permissions
- Use `--log-level debug` for detailed diagnostic information

**"Permission denied" errors**
- Verify you have `resourcemanager.folders.delete` and `resourcemanager.projects.delete` permissions
- Check if folders/projects have any protection policies
- Use `--log-level trace` to see exact gcloud commands being executed

**"Folder not empty" errors**
- This shouldn't happen with proper post-order traversal
- Check logs for failed deletions of child resources with `--log-level debug`
- Some resources may require manual cleanup (e.g., billing accounts, liens)

**API Rate Limiting with Concurrency**
- Reduce `--concurrency-limit` value (try 3-5 for moderate hierarchies)
- Add delays between operations by using sequential mode temporarily
- Check GCP quotas and limits for your project
- Use `--log-level debug` to see timing of API calls

**High Memory Usage with Concurrency**
- Reduce `--concurrency-limit` value
- Process smaller subsections of the hierarchy
- Monitor system resources and adjust accordingly

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Test with different concurrency settings: `--concurrency --concurrency-limit 3` and `--concurrency --concurrency-limit 10`
6. Test with different log levels: `--log-level debug` and `--log-level trace`
7. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Important Notes

⚠️ **Resource Deletion Behavior:**
- **Projects**: Enter a 30-day pending deletion period and can be recovered via the [Cloud Resource Manager](https://console.cloud.google.com/cloud-resource-manager?pendingDeletion=true)
- **Folders**: Are permanently deleted and **cannot be recovered**
- Always use `--dry-run` first and review the tree visualization before proceeding with actual deletion
- Use appropriate log levels (`--log-level debug`) when investigating issues or validating complex hierarchies
- Start with lower concurrency limits and increase gradually based on your hierarchy size and network conditions

⚠️ **Concurrency Considerations:**
- **Start conservatively**: Begin with `--concurrency-limit 3-5` and increase based on performance
- **Monitor API limits**: GCP has rate limits that may be hit with high concurrency
- **Test thoroughly**: Always test concurrent operations with `--dry-run` before actual deletion
- **Debug with lower concurrency**: Use `--concurrency-limit 1-3` with `--log-level debug` for troubleshooting
