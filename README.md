# Distributed Script Management System

A Go-based distributed system for executing scripts and commands across multiple agents running in Docker containers. This project provides centralized management of system monitoring, security scanning, and automation tasks across a distributed environment.

## Overview

This system consists of a central script manager that can send bash scripts to multiple agents running in separate Docker containers. Each agent executes the received scripts and returns the results to the central manager, which aggregates and displays the output.

## Architecture

### Core Components

**Script Manager (Go)**
- Central controller that reads script files and distributes them to agents
- Handles concurrent communication with multiple agents
- Aggregates and displays results from all agents
- Supports both host and container-specific scripts

**Agent System (Go)**
- TCP-based communication server running on each agent
- Executes received scripts in the local environment
- Handles multi-line scripts by creating temporary files
- Returns execution results to the script manager

**Docker Integration**
- Agents run in isolated Ubuntu containers
- Each container has its own network namespace and filesystem
- Scripts execute within container boundaries
- System information reflects container environment

### Communication Flow

1. Script Manager reads a bash script file
2. Script content is sent via TCP to all registered agents
3. Each agent receives the script and executes it locally
4. Execution results are returned to the script manager
5. Script manager aggregates and displays all results

## Features

### Script Categories

**Host Scripts** (for physical machine monitoring)
- `system_monitor.sh`: Hardware information, temperature, network traffic
- `security_scan.sh`: Port scanning, malicious software detection, user activity
- `system_info.sh`: Basic system information and status

**Container Scripts** (for Docker container monitoring)
- `container_monitor.sh`: Container-specific system information
- `container_security.sh`: Container security analysis
- `backup_files.sh`: File backup operations
- `cleanup_logs.sh`: Log file management
- `security_check.sh`: Basic security verification

### System Capabilities

**Distributed Execution**
- Execute scripts on multiple agents simultaneously
- Real-time result aggregation
- Concurrent processing with Go goroutines

**System Monitoring**
- CPU, RAM, and disk usage monitoring
- Network interface and traffic analysis
- Process monitoring and resource tracking

**Security Features**
- Port scanning and network analysis
- User activity and sudo command logging
- File system integrity checks
- Container-specific security scanning

**Automation Support**
- File backup and cleanup operations
- Log management and rotation
- System maintenance tasks
- Custom script execution

## Installation and Setup

### Prerequisites

- Go 1.19 or higher
- Docker
- Ubuntu/Debian system (for container agents)
- Basic networking knowledge

### Building the System

```bash
# Build the agent
cd agent
go build -o agent agent.go

# Build the script manager
cd script-manager
go build -o script_manager script_manager.go
```

### Setting Up Agents

```bash
# Create Ubuntu containers for agents
docker run -d --name ubuntu-agent1 -p 9001:9001 ubuntu:latest tail -f /dev/null
docker run -d --name ubuntu-agent2 -p 9002:9002 ubuntu:latest tail -f /dev/null
docker run -d --name ubuntu-agent3 -p 9003:9003 ubuntu:latest tail -f /dev/null

# Install required packages in containers
docker exec ubuntu-agent1 apt update && docker exec ubuntu-agent1 apt install -y procps iproute2 net-tools lsof coreutils
docker exec ubuntu-agent2 apt update && docker exec ubuntu-agent2 apt install -y procps iproute2 net-tools lsof coreutils
docker exec ubuntu-agent3 apt update && docker exec ubuntu-agent3 apt install -y procps iproute2 net-tools lsof coreutils

# Copy agent binary to containers
docker cp agent/agent ubuntu-agent1:/agent
docker cp agent/agent ubuntu-agent2:/agent
docker cp agent/agent ubuntu-agent3:/agent

# Copy scripts to containers
docker cp scripts/container/. ubuntu-agent1:/scripts
docker cp scripts/container/. ubuntu-agent2:/scripts
docker cp scripts/container/. ubuntu-agent3:/scripts

# Start agents in containers
docker exec -d ubuntu-agent1 /agent 9001
docker exec -d ubuntu-agent2 /agent 9002
docker exec -d ubuntu-agent3 /agent 9003
```

## Usage

### Running the Script Manager

```bash
# Start the script manager
./run_script_manager.sh
```

### Available Scripts

The script manager will display available scripts:

**Host Scripts (for physical machine)**
- `scripts/host/system_monitor.sh`
- `scripts/host/security_scan.sh`
- `scripts/host/system_info.sh`

**Container Scripts (for Docker containers)**
- `scripts/container/container_monitor.sh`
- `scripts/container/container_security.sh`
- `scripts/container/backup_files.sh`
- `scripts/container/cleanup_logs.sh`
- `scripts/container/security_check.sh`

### Example Usage

```bash
# Execute container monitoring script
scripts/container/container_monitor.sh

# Execute security scan on all agents
scripts/container/container_security.sh

# Run system monitoring on host
scripts/host/system_monitor.sh
```

## Technical Details

### Network Configuration

- Script Manager: Runs on host machine
- Agent 1: Container ubuntu-agent1, port 9001
- Agent 2: Container ubuntu-agent2, port 9002
- Agent 3: Container ubuntu-agent3, port 9003

### Script Execution Process

1. Script Manager reads script file content
2. Content is sent via TCP to all agents
3. Each agent creates temporary script file
4. Script is executed with bash
5. Output is captured and returned
6. Results are aggregated and displayed

### Error Handling

- Connection failures are reported per agent
- Script execution errors are captured
- Timeout handling for long-running scripts
- Graceful degradation when agents are unavailable

## Development

### Project Structure

```
bash-king/
├── agent/                 # Agent source code
│   └── agent.go         # TCP server and script executor
├── script-manager/       # Script manager source code
│   └── script_manager.go # Central controller
├── scripts/              # Bash scripts
│   ├── host/            # Host-specific scripts
│   └── container/       # Container-specific scripts
├── monitoring/           # Monitoring agent (alternative)
├── server/              # Legacy server implementation
├── old_versions/        # Previous implementations
└── run_script_manager.sh # Convenience script
```

### Adding New Scripts

1. Create script in appropriate directory (`host/` or `container/`)
2. Make script executable: `chmod +x scripts/container/new_script.sh`
3. Script will be available in script manager

### Extending the System

- Add new agent types by implementing the TCP protocol
- Create custom scripts for specific use cases
- Extend monitoring capabilities with additional metrics
- Implement authentication and security features

## Use Cases

**System Administration**
- Monitor multiple servers simultaneously
- Execute maintenance scripts across infrastructure
- Collect system information from distributed environments

**DevOps Operations**
- Deploy configuration changes to multiple containers
- Run health checks across microservices
- Execute backup and cleanup operations

**Security Operations**
- Scan multiple systems for vulnerabilities
- Monitor network activity across containers
- Audit user activity and system changes

**Development and Testing**
- Test scripts in isolated container environments
- Validate system configurations across platforms
- Debug distributed system issues

## Limitations and Considerations

**Container Limitations**
- Some hardware information may not be available in containers
- Temperature sensors and certain system calls may be restricted
- Network information reflects container networking

**Security Considerations**
- Scripts execute with container privileges
- No authentication implemented in current version
- Consider network security for production deployment

**Performance Notes**
- Script execution time depends on container resources
- Large script outputs may impact network performance
- Concurrent execution may stress system resources

## Future Enhancements

**Planned Features**
- Web-based dashboard for script management
- Authentication and authorization system
- Real-time monitoring and alerting
- Script scheduling and automation
- File transfer capabilities between agents

**Potential Improvements**
- Support for different container platforms
- Integration with monitoring systems
- Advanced logging and audit trails
- Plugin architecture for custom scripts
- High availability and failover support

## Troubleshooting

**Common Issues**

*Agents not responding*
- Check if containers are running: `docker ps`
- Verify agent processes: `docker exec container_name ps aux | grep agent`
- Check network connectivity: `telnet localhost 9001`

*Scripts not executing*
- Verify script permissions: `ls -la scripts/`
- Check script syntax: `bash -n script.sh`
- Review agent logs for errors

*Empty output from containers*
- Ensure required packages are installed in containers
- Check if commands are available in container environment
- Verify script compatibility with container OS

**Debug Information**

Enable debug mode by checking agent logs:
```bash
docker logs ubuntu-agent1
docker logs ubuntu-agent2
docker logs ubuntu-agent3
```

## License

This project is open source and available under the MIT License. See LICENSE file for details.

## Contributing

Contributions are welcome. Please read the contributing guidelines before submitting pull requests.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For issues and questions, please open an issue on the project repository or contact the maintainers. 