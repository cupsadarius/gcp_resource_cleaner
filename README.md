### GCP Resource Cleaner

This is a simple tool that can be used to clean up GCP resources. It can be used to delete unused resources, or to clean up resources that are no longer needed.

### Features

- Delete unused resources: This tool can be used to delete unused resources, such as folders and projects, that are no longer needed.

### Installation

To install the tool, you can use the following command:

```
go install github.com/cupsadarius/gcp_resource_cleaner@latest
```

### Usage

To use the tool, you need to have the following:

- Gcloud CLI installed and configured with a user that has the necessary permissions to delete resources (folders and projects).
- A GCP project ID as the root of the resource hierarchy.

Once you have these, you can run the tool with the following command:

```
gcp_resource_cleaner check-health
```

This will check the health of the GCP resources and print out any issues that need to be addressed.

You can also run the tool with the following command to delete unused resources:

```
gcp_resource_cleaner delete --folder-id=<folder-id> --dry-run <dry-run>
```

This will delete unused resources and print out the resources that will be deleted. You can then review the list and decide whether to proceed with the deletion.

If you want to delete resources without a dry run, you can use the following command:

```
gcp_resource_cleaner delete --folder-id=<folder-id>
```

This will delete unused resources without a dry run.

