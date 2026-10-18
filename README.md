# CSV Validator

A Go command-line tool that checks required columns, missing required values, and duplicate IDs. Uses only the Go standard library.

## Build and run

Requires Go 1.27.2 or newer. Run commands from the project root:

```bash
go build -o bin/csv-validator ./cmd/csv-validator
./bin/csv-validator --file users.csv --required id,name,email --id id
```

All three flags are required:

| Flag | Meaning |
| --- | --- |
| `--file` | Path to a CSV file |
| `--required` | Comma-separated required column names, without added spaces |
| `--id` | ID column name; must also appear in `--required` |

Show help:

```bash
./bin/csv-validator --help
```

Use the built binary when checking exit codes. `go run` does not preserve the program's validation exit code directly.

## Example

Save as `users.csv`:

```csv
id,name,email
42,Ada,ada@example.com
 42 ,Grace,
```

Run:

```bash
./bin/csv-validator --file users.csv --required id,name,email --id id
```

Output:

```text
record 3, column "email": required value missing
record 3, column "id": duplicate ID
```

Valid input prints `CSV valid`.

### Validation rules

- Headers and configured column names match exactly, including case and surrounding whitespace.
- Empty or whitespace-only required values are missing. Optional values may be blank.
- IDs are compared after trimming surrounding whitespace. `42` and ` 42 ` are duplicates; `ABC` and `abc` are distinct.
- Blank IDs produce missing-value findings, not duplicate-ID findings.
- Blank or duplicate headers produce findings and stop data-record validation. Missing required headers also stop data-record validation.
- Empty files, malformed CSV, and records with the wrong field count produce errors.
- Quoted commas and multiline fields follow Go's `encoding/csv` rules. Empty physical lines are ignored by the CSV reader.
- Record numbers include the header as record 1; multiline fields do not create extra records.
- Records are read incrementally. Unique IDs and accumulated findings remain in memory.

### Exit codes

| Code | Meaning |
| --- | --- |
| `0` | CSV valid, or help requested |
| `1` | Validation findings |
| `2` | Invalid arguments/configuration, file access, parsing, or report-writing error |

Findings and success messages go to stdout. Usage and processing errors go to stderr. Findings collected before a parsing failure are printed, but the final exit code is `2` because validation could not finish.

## Docker

Build locally:

```bash
docker build -t csv-validator:local .
```

Run against a file in the current directory:

```bash
docker run --rm --network=none --read-only \
  --mount "type=bind,source=$(pwd),target=/data,readonly" \
  csv-validator:local \
  --file /data/users.csv --required id,name,email --id id
```

The image contains a static binary, runs as non-root UID/GID `65532:65532`, and uses `/data` as its working directory. Mounted input files must be readable and their directories traversable by that user. Input is mounted read-only; no ports or network access are needed.

After publishing, replace `csv-validator:local` with `YOUR_DOCKERHUB_USERNAME/csv-validator:YOUR_RELEASE_TAG`.

## Tests

```bash
go test ./... -count=1
go vet ./...
```

## Publish to Docker Hub through GitHub Releases

Workflow: [`.github/workflows/docker-publish.yml`](.github/workflows/docker-publish.yml).

1. Create a Docker Hub repository named `csv-validator` under your account.
2. Create a Docker Hub access token with Read & Write permission for publishing.
3. In the GitHub repository, open **Settings → Secrets and variables → Actions** and add repository secrets:
   - `DOCKERHUB_USERNAME`: Docker Hub username.
   - `DOCKERHUB_TOKEN`: Docker Hub access token. Never commit this token.
4. Commit and push the project, including the workflow. Ensure GitHub Actions is enabled.
5. Create and publish a GitHub Release, for example with tag `v1.0.0`. The tagged commit must contain the workflow.
6. Check the **Actions** tab for the publishing result.

The workflow tests and vets the code before building and pushing a Linux AMD64 image:

- Every published release gets its release tag, such as `YOUR_DOCKERHUB_USERNAME/csv-validator:v1.0.0`.
- Releases not marked as prereleases also update `YOUR_DOCKERHUB_USERNAME/csv-validator:latest`.
- Prereleases do not update `latest`.
- Draft releases do not publish images until published.

Publishing an older stable release also moves `latest` to that release. Use explicit version tags when reproducibility matters. ARM64 images are not currently published.
