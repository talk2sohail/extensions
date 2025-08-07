package main

import (
	"context"
	"flag"
	"log"
	"strings"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/table"
)

// main is the entry point for the osquery extension.
func main() {
	// osquery-go requires a socket path to be specified.
	socket := flag.String("socket", "", "Path to the osquery extensions UNIX domain socket")
	flag.Parse()
	if *socket == "" {
		log.Fatalf("Usage: container_processes_plugin --socket SOCKET_PATH")
	}

	// Create an extension server.
	server, err := osquery.NewExtensionManagerServer("container_processes_extension", *socket)
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	// Create and register the container_processes table plugin.
	server.RegisterPlugin(table.NewPlugin("container_processes", ContainerProcessesColumns(), ContainerProcessesGenerate))

	// Start the server and wait for it to be shut down.
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

// ContainerProcessesColumns defines the schema for our table.
func ContainerProcessesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		// "host_pid" is the process ID on the host OS running the Docker daemon (e.g., a Linux VM on Windows/macOS), not inside the container's PID namespace.
		table.IntegerColumn("host_pid"),
		table.TextColumn("name"),
		table.TextColumn("container_id"),
		table.TextColumn("container_name"),
		table.TextColumn("container_image"),
	}
}

// ContainerProcessesGenerate is the function that osquery calls to generate the table rows.
func ContainerProcessesGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	// Create a new Docker client. It will use the DOCKER_HOST environment variable
	// or the default socket path.
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Error creating Docker client: %v. Is Docker running?", err)
		// If we can't connect to Docker, we return an empty table.
		return []map[string]string{}, nil
	}

	// List all running containers.
	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{All: true})
	if err != nil {
		log.Printf("Error listing containers: %v", err)
		return nil, err
	}

	var results []map[string]string

	// Iterate over each container to get its processes.
	for _, container := range containers {
		// The ContainerTop command is equivalent to `docker top`.
		processes, err := cli.ContainerTop(ctx, container.ID, nil)
		if err != nil {
			log.Printf("Error getting processes for container %s: %v", container.ID, err)
			continue // Move to the next container
		}

		// Find the column indexes for PID and COMMAND from the Titles.
		pidIndex := findColumnIndex(processes.Titles, "PID")
		// The command column can be either "CMD" or "COMMAND".
		cmdIndex := findColumnIndex(processes.Titles, "COMMAND")
		if cmdIndex == -1 {
			cmdIndex = findColumnIndex(processes.Titles, "CMD")
		}

		if pidIndex == -1 || cmdIndex == -1 {
			log.Printf("Could not find PID or COMMAND/CMD columns for container %s. Titles: %v", container.ID, processes.Titles)
			continue
		}

		// Get a clean container name.
		containerName := "unknown"
		if len(container.Names) > 0 {
			// Names are often prefixed with a forward slash.
			containerName = strings.TrimPrefix(container.Names[0], "/")
		}

		// Create a row for each process found in the container.
		for _, proc := range processes.Processes {
			if len(proc) <= pidIndex || len(proc) <= cmdIndex {
				continue
			}
			row := map[string]string{
				"host_pid": proc[pidIndex],
				"name":     proc[cmdIndex],
				// Use the short container ID for readability.
				"container_id":    container.ID[:12],
				"container_name":  containerName,
				"container_image": container.Image,
			}
			results = append(results, row)
		}
	}

	return results, nil
}

// findColumnIndex is a helper to find the index of a column in a slice of strings.
// It returns -1 if the column is not found.
func findColumnIndex(slice []string, value string) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}
