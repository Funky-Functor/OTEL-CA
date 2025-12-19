# OTEL-CA
Test tool to check OpenTelemetry configurations.

## How to transpile the dingo file?
1. Install Dingo if you haven't already. You can find it at [Dingo's GitHub repository](https://github.com/MadAppGang/dingo/tree/main#installation)
2. Run the following command in your terminal:
   ```bash
   dingo go main.dingo
   ```
## Docker

### Build the Docker image
```bash
docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/funkyfunctor-chris/otel_ca:202512182135 .
```

:warning: This will build the image for multiple architectures (amd64 and arm64).

### Run the Docker container
```bash
docker run -v $(pwd)/config.json:/app/config.json ghcr.io/funkyfunctor-chris/otel_ca:202512182135 -config /app/config.json
```

### Push the Docker image to GitHub Container Registry
You first need to login to GitHub Container Registry:
```bash
echo <Personal Access Token> | docker login ghcr.io -u FunkyFunctor-Chris --password-stdin
```

Once you are logged in, you can push the image:
```bash
docker push ghcr.io/funkyfunctor-chris/otel_ca:202512182135
```

Once it has been pushed, the image should be visible at [the following location](https://github.com/users/FunkyFunctor-Chris/packages/container/package/otel_ca)

## DQL Query to find the generated traces

:warning: Practically, you would want to adjust the value of your test marker to make sure you find the request you are looking for.

```
fetch spans, from: -30m, samplingRatio: 1, scanLimitGBytes: 500

// based on the applied filters
| filter (isNull(dt.agent.module.id) AND (isNull(telemetry.exporter.name) or not(matchesValue(telemetry.exporter.name, "odin"))) AND not matchesValue(dt.openpipeline.source, "oneagent")) /* span.source condition is based on these attributes */ 
	 and matchesValue(`test.marker`, "test")

| limit 1000

// construct fields
| fieldsAdd request.status_code = if(request.is_failed, "Failure", else: "Success")
| fieldsAdd span.source = if((isNotNull(dt.agent.module.id) or matchesValue(telemetry.exporter.name, "odin") or matchesValue(telemetry.sdk.name, "oneagent") or matchesValue(dt.openpipeline.source, "oneagent")), "OneAgent", else: "OpenTelemetry")

// prepare fields
| fieldsAdd http.response.status_code = coalesce(http.response.status_code, toLong(http.status_code))
| fieldsAdd k8s.workload.name = coalesce(k8s.workload.name, dt.kubernetes.workload.name)

// always limit fields for performance
| fields
    dt.entity.process_group,
    dt.entity.service,
    duration,
    endpoint.name,
    gen_ai.completion.0.content,
    gen_ai.prompt.0.content,
    http.response.status_code,
    k8s.workload.name,
    dt.entity.cloud_application,
    request.status_code,
    span.source,
    start_time,
    dt.agent.module.id,
    span.id,
    trace.id,
    dt.system.sampling_ratio
// add entity lookups
| fieldsAdd dt.entity.process_group.entity.name = entityAttr(dt.entity.process_group, "entity.name"), dt.entity.service.entity.name = entityAttr(dt.entity.service, "entity.name")
```