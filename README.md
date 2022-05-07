# nomad-playground

A CLI tool to create/delete a basic job for running Docker container within [Hashicorp Nomad](https://www.nomadproject.io/) orchestrator.

## Troubleshooting port mapping for Docker containers

Run this Bash command:
```bash
for port in $(docker ps --format "{{.Ports}}"); do echo $port; done | grep tcp | cut -d':' -f 2
```
