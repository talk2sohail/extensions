package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/table"
)

var (
	// socket is the path to the osquery extensions UNIX domain socket.
	socket = flag.String("socket", "", "Path to the osquery extensions UNIX domain socket")
	// containerIDLength is the length of the container ID to display.
	containerIDLength = flag.Int("container_id_length", 12, "Length of the container ID to display")
)

// Plugin is the main struct for our osquery extension.
// It holds a reference to the Docker client.
type Plugin struct {
	cli *client.Client
}

// NewPlugin creates a new instance of our plugin.
func NewPlugin() (*Plugin, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %w", err)
	}
	return &Plugin{cli: cli}, nil
}

// main is the entry point for the osquery extension.
func main() {
	flag.Parse()
	if *socket == "" {
		log.Fatalf("Usage: container_processes_plugin --socket SOCKET_PATH")
	}

	server, err := osquery.NewExtensionManagerServer("container_processes_extension", *socket)
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	plugin, err := NewPlugin()
	if err != nil {
		log.Fatalf("Error creating plugin: %s\n", err)
	}

	server.RegisterPlugin(table.NewPlugin("container_processes", ContainerProcessesColumns(), plugin.ContainerProcessesGenerate))

	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

// ContainerProcessesColumns defines the schema for our table.
func ContainerProcessesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.IntegerColumn("host_pid"),
		table.TextColumn("name"),
		table.TextColumn("container_id"),
		table.TextColumn("container_name"),
		table.TextColumn("container_image"),
	}
}

// ContainerProcessesGenerate is the function that osquery calls to generate the table rows.
func (p *Plugin) ContainerProcessesGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	containers, err := p.cli.ContainerList(ctx, containertypes.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %w", err)
	}

	var results []map[string]string

	for _, container := range containers {
		processes, err := p.cli.ContainerTop(ctx, container.ID, nil)
		if err != nil {
			log.Printf("Error getting processes for container %s: %v", container.ID, err)
			continue
		}

		pidIndex := findColumnIndex(processes.Titles, "PID")
		cmdIndex := findColumnIndex(processes.Titles, "COMMAND")
		if cmdIndex == -1 {
			cmdIndex = findColumnIndex(processes.Titles, "CMD")
		}

		if pidIndex == -1 || cmdIndex == -1 {
			log.Printf("Could not find PID or COMMAND/CMD columns for container %s. Titles: %v", container.ID, processes.Titles)
			continue
		}

		containerName := "unknown"
		if len(container.Names) > 0 {
			containerName = strings.TrimPrefix(container.Names[0], "/")
		}

		for _, proc := range processes.Processes {
			if len(proc) <= pidIndex || len(proc) <= cmdIndex {
				continue
			}
			row := map[string]string{
				"host_pid":        proc[pidIndex],
				"name":            proc[cmdIndex],
				"container_id":    container.ID[:*containerIDLength],
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
