# Contributing

Thanks for helping improve the public bug bounty program dataset.

## Scope

This repository tracks public bug bounty and responsible disclosure programs used by ProjectDiscovery's Chaos dataset.

Most contributions change:

- `src/data.yaml`
- `src/data.schema.json` (only when the data model changes)

Generated output:

- `dist/data.json` is generated from `src/data.yaml`

## Requirements

Every program entry in `src/data.yaml` must include:

- `name` as a string
- `url` as an `http` or `https` URL
- `bounty` as a boolean
- `domains` as a list of root/apex domains

Domain rules:

- Use only domain names in `domains`
- Do not include wildcards like `*.example.com`
- Do not include full URLs like `https://example.com`
- Prefer primary/apex domains (for example: `example.com`, `example.co.uk`)
- Keep each program's domain list unique

## Local Validation

Run checks before opening a pull request.

1. Build generated JSON from YAML:

```bash
make compile
```

2. Validate generated data against schema:

```bash
make test
```

3. Check for duplicate domains across all programs:

```bash
make duplicate-domains
```

4. Validate domain formatting rules:

```bash
make validate-domains
```

Optional URL policy checks:

```bash
make policy-checks
```

## Editing Guidelines

- Keep entries sorted only if you are intentionally doing a full reorder; otherwise, minimize unrelated movement.
- Keep changes focused. Avoid mixing schema refactors and large data updates in one PR.
- If a domain is intentionally removed or changed, add context in the PR description.

## Pull Request Guidelines

Include the following in your PR:

- A short summary of what was added, removed, or corrected
- Why the change is needed (source link, program page, or verification details)
- Confirmation that local checks passed

PRs that fail CI checks (`compile`, `test`, duplicate-domain check, domain validation) will need to be fixed before merge.

## Discussions and Questions

- Use GitHub Discussions for ideas and broader suggestions.
- Use Issues for concrete bugs, data problems, or reproducible validation failures.
