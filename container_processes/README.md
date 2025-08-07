## container_processes osquery Extension

This extension adds a custom `container_processes` table to osquery, allowing you to view and correlate running processes inside Docker containers with their host context. It provides visibility into containerized workloads by exposing process details, container metadata, and image information, making it easier to monitor, audit, and secure your container environment.
### What does it do?

- Lists all running processes inside each Docker container.
- Shows the host PID, process name, container ID, container name, and container image for each process.
- Enables correlation between host processes and their container context for security and monitoring.

---

## Usage Instructions (Linux)

### Prerequisites
- Go (for building the extension)
- Docker daemon running and accessible
- osquery installed (https://osquery.io/downloads/official)

### Build the Extension

```sh
cd extensions/container_processes
go build -o container_processes_plugin main.go
```

### Run the Extension with osqueryd

1. Start osqueryd with the extension:
   ```sh
   osqueryd --extension /path/to/container_processes_plugin --flagfile /etc/osquery/osquery.flags
   ```
   Replace `/path/to/container_processes_plugin` with the actual path to your built binary.

2. Query the table using osqueryi or osqueryd:
   ```sql
   SELECT * FROM container_processes;
   ```

### Notes
- The extension must be able to connect to the Docker daemon (usually via `/var/run/docker.sock`).
- You may need to run osquery as root or with appropriate permissions to access Docker and process information.

---

For more details, see the source code in `main.go`.
