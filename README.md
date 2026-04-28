# grpcurl-batch

A wrapper around [grpcurl](https://github.com/fullstorydev/grpcurl) for running batched gRPC requests from YAML definitions with retries and output formatting.

---

## Installation

```bash
go install github.com/youruser/grpcurl-batch@latest
```

> **Prerequisite:** [grpcurl](https://github.com/fullstorydev/grpcurl) must be installed and available on your `PATH`.

---

## Usage

Define your requests in a YAML file:

```yaml
# requests.yaml
host: localhost:50051
requests:
  - method: mypackage.MyService/GetUser
    data: '{"id": "123"}'
    retries: 3
  - method: mypackage.MyService/ListItems
    data: '{"page": 1}'
    retries: 1
    timeout: 5s
```

Run the batch:

```bash
grpcurl-batch --file requests.yaml --output json
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the YAML request definition file | `requests.yaml` |
| `--output` | Output format: `json`, `pretty`, `summary` | `pretty` |
| `--concurrency` | Number of concurrent requests | `1` |
| `--retry-delay` | Delay between retries | `1s` |

### Example Output

```
[1/2] mypackage.MyService/GetUser ... OK (142ms)
[2/2] mypackage.MyService/ListItems ... FAILED (attempt 1/1)

Summary: 1 succeeded, 1 failed
```

---

## License

[MIT](LICENSE)