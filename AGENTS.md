# Repository rules

## What this is

A community port of [Lody](https://github.com/lbr77/lody) — the desktop client
at `/Users/salt/Code/lody` — to [MyGo](https://github.com/egoist/mygo), the
Go UI framework at `/Users/salt/Code/mygo`. Where Lody draws its interface
in a web view, this client draws on the GPU with MyGo's `ui` package, no DOM.

The port follows Lody's desktop app one to one: layout, design tokens, wording,
motion, and the protocol it speaks to Lody Cloud. This fork strips the features
its owner does not use: AI usage, billing and referrals, publishing sessions
to share.lody.ai, and share-as-image.

## Reference codebases

When writing or changing code, consult the **latest** source of both:

- `/Users/salt/Code/lody` — the upstream Lody desktop app (TypeScript,
  Electron). Its `packages/components/src/` holds the React UI this port
  mirrors; `apps/cli/` and `packages/shared/` the protocol, session and
  catalog semantics. Match its behavior before inventing your own.
- `/Users/salt/Code/mygo` — the framework. Its `ui` package supplies the
  elements, motion and input this app draws with; `cmd/mygo` the `dev`,
  `build` and `generate` tooling. Check its docs and source for the API the
  vendored `github.com/egoist/mygo` version exposes.

Both live outside this repository; read them directly rather than relying on
memory of an older version. The vendored MyGo version is in `go.mod`; the
module cache copy at `$(go env GOMODCACHE)/github.com/egoist/mygo@<version>`
shows exactly what this build links.

## Layout

- `main.go` — entry: app name, single instance, the window.
- `internal/app` — the whole client: views, settings, sidebar, composer,
  session lifecycle. Shared controls are in `controls.go`.
- `internal/catalog` — the workspace and machine catalogs merged from the
  Flock streams; projects, sessions, agents.
- `internal/cloud`, `internal/streams`, `internal/flock` — Lody Cloud API,
  Loro/SSE transport, the Flock document model.
- `locales` — `en.json` and `zh-Hans.json`; `go tool mygo generate` writes
  `catalog_gen.go` after they change.

## Build and test

```sh
CGO_ENABLED=1 go run .          # run
CGO_ENABLED=1 go test ./...     # tests, headless UI tests included
```

CGO stays on: `loro-go` links a prebuilt Rust library through it.
