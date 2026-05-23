# gh-app-inspector

Inspect a GitHub App installation from the command line: App metadata, installation details, granted permissions, accessible repositories, and core rate limit — all rendered as rounded tables.

## Install

```sh
go install github.com/wozniakpl/gh-app-inspector@latest
```

## Usage

```sh
gh-app-inspector \
  --app-id 123456 \
  --installation-id 78901234 \
  --pem ~/path/to/key.pem
```

Or via environment variables:

```sh
export GH_APP_ID=123456
export GH_INSTALLATION_ID=78901234
export GH_APP_PEM=~/path/to/key.pem
gh-app-inspector
```

The PEM must be the App's private key (`-----BEGIN RSA PRIVATE KEY-----`), downloaded from the GitHub App settings page.

## Example output

```
╭───────────────────────────────────────────────────╮
│ GitHub App                                        │
├──────────┬────────────────────────────────────────┤
│ Name     │ example-app                            │
│ Slug     │ example-app                            │
│ ID       │ 123456                                 │
│ Owner    │ example-org                            │
│ HTML URL │ https://github.com/apps/example-app    │
│ Created  │ 2025-01-01T00:00:00Z                   │
╰──────────┴────────────────────────────────────────╯

╭─────────────────────────────────────────────╮
│ Installation                                │
├──────────────────────┬──────────────────────┤
│ ID                   │ 78901234             │
│ Account              │ example-org          │
│ Account type         │ Organization         │
│ Target type          │ Organization         │
│ Repository selection │ selected             │
│ Events               │ [push pull_request]  │
│ Created              │ 2025-01-01T00:00:00Z │
│ Updated              │ 2025-06-01T00:00:00Z │
╰──────────────────────┴──────────────────────╯

╭─────────────────────────╮
│ Permissions (4)         │
├────────────────┬────────┤
│ SCOPE          │ ACCESS │
├────────────────┼────────┤
│ contents       │ read   │
│ issues         │ write  │
│ metadata       │ read   │
│ pull_requests  │ write  │
╰────────────────┴────────╯

╭───────────────────────────────────────────────────────────────────────────╮
│ Accessible repositories (3)                                               │
├─────────────────────────┬─────────┬────────────────┬──────────────────────┤
│ REPOSITORY              │ PRIVATE │ DEFAULT BRANCH │ PUSHED AT            │
├─────────────────────────┼─────────┼────────────────┼──────────────────────┤
│ example-org/alpha       │ false   │ main           │ 2025-05-01T12:00:00Z │
│ example-org/beta        │ true    │ main           │ 2025-05-02T12:00:00Z │
│ example-org/gamma       │ true    │ main           │ 2025-05-03T12:00:00Z │
╰─────────────────────────┴─────────┴────────────────┴──────────────────────╯

╭──────────────────────────────────╮
│ Rate limit (core)                │
├───────────┬──────────────────────┤
│ Limit     │ 5000                 │
│ Remaining │ 4987                 │
│ Resets    │ 2025-06-01T13:00:00Z │
╰───────────┴──────────────────────╯
```

## Releases

Versioning is automated by [tag-it](https://github.com/wozniakpl/tag-it). PR titles follow Conventional Commits; merges to `main` cut a new semver tag automatically.
