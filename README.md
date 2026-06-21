# jsonstat

A command-line utility written in Go for processing JSON Lines (JSONL) files.

The application reads a stream of JSON objects, filters records by a specified field, outputs matching records, and provides processing statistics.

## Features

- Read input from a file or stdin
- Filter records by:
  - `id`
  - `name`
  - `role`
  - `salary`
- Limit the number of processed records
- Layered configuration support:
  - configuration file
  - environment variables
  - command-line flags
- JSON Lines (JSONL) support
- Unit tests, fuzz tests, and benchmark tests

---

## Technical Highlights

- Layered configuration loading:
  `defaults → config file → environment variables → CLI flags`
- Stream processing using `io.Reader` and `io.Writer`
- Table-driven tests
- Fuzz testing
- Benchmark testing
- Error wrapping using `%w`

---

## Installation

Clone the repository:

```bash
git clone https://github.com/kihcnxhelp/jsonstat.git

cd jsonstat
```

Build the application:

```bash
go build -o jsonstat .
```

Or run directly:

```bash
go run .
```

---

## Input Format

The application expects input in JSON Lines format, where each line contains a separate JSON object:

```json
{"id":1,"name":"Alice","role":"admin","salary":100}
{"id":2,"name":"Bob","role":"user","salary":50}
{"id":3,"name":"Carol","role":"admin","salary":120}
```

---

## Usage

Print all records:

```bash
jsonstat -input employees.jsonl
```

Filter by role:

```bash
jsonstat \
  -input employees.jsonl \
  -field role \
  -value admin
```

Limit processed records:

```bash
jsonstat \
  -input employees.jsonl \
  -max 100
```

Read from stdin:

```bash
cat employees.jsonl | jsonstat -input -
```

---

## Configuration

Configuration values are applied in the following order:

```text
defaults
→ config file
→ environment variables
→ command-line flags
```

Each layer overrides values from the previous one.

### Configuration File

By default, the application looks for:

```text
config.json
```

Example:

```json
{
  "log_level": "info",
  "max_records": 100
}
```

Specify a custom configuration file:

```bash
export JSONSTAT_CONFIG=./configs/dev.json
```

---

### Environment Variables

| Variable | Description |
|-----------|-------------|
| JSONSTAT_INPUT_FILE | Input file path |
| JSONSTAT_FILTER_FIELD | Filter field |
| JSONSTAT_FILTER_VALUE | Filter value |
| JSONSTAT_MAX_RECORDS | Maximum number of records |
| JSONSTAT_LOG_LEVEL | Log level |
| JSONSTAT_CONFIG | Path to configuration file |

---

## Example

Filter administrators:

```bash
jsonstat \
  -input employees.jsonl \
  -field role \
  -value admin
```

Output:

```json
{"id":1,"name":"Alice","role":"admin","salary":100}
{"id":3,"name":"Carol","role":"admin","salary":120}
```

Statistics:

```text
total=3 matched=2 skipped=1
```

---

## Testing

Run unit tests:

```bash
go test ./...
```

Run fuzz tests:

```bash
go test -fuzz=Fuzz ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

---

## Project Structure

```text
internal/
├── config
│   └── configuration loading and validation
│
└── processor
    └── JSONL stream processing
```

The `processor` package operates through `io.Reader` and `io.Writer` interfaces, making the business logic independent from files, stdin/stdout, and easier to test.

---

## Future Improvements

Potential enhancements:

- Structured logging with `log/slog`
- Additional filter operators (`>`, `<`, `contains`, etc.)
- CSV export
- Output formatting options
- Parallel processing for large datasets
- Docker support
- CI/CD pipeline with GitHub Actions

---

## License

MIT