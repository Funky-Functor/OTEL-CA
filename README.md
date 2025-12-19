# OTEL-CA
Test tool to check OpenTelemetry configurations.

## How to transpile the dingo file?
1. Install Dingo if you haven't already. You can find it at [Dingo's GitHub repository](https://github.com/MadAppGang/dingo/tree/main#installation)
2. Run the following command in your terminal:
   ```bash
   dingo build main.dingo
   ```
## Docker

### Build the Docker image
```bash
docker build -t otel-ca .
```

### Run the Docker container
```bash
docker run -v $(pwd)/config.json:/app/config.json otel_ca:latest -config /app/config.json
```
