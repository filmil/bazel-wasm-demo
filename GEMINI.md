# Bazel Go WASM Demo - Gemini Instructions

This project is a Bazel-managed Go WASM application. Adhere to the following
mandates and workflows.

## Build System & Toolchain

- **Build System:** Bazel (version `9.0.1`).
- **Primary Language:** Go.
- **Commands:**
  - **NEVER** run `go` directly. Always use `bazel run @rules_go//go --
    <args>`.
  - Use `bazel run //:gazelle` to update build rules.
  - Use `bazel mod tidy` to update `MODULE.bazel`.
  - Build: `bazel build //...`
  - Test: `bazel test //...`

## Engineering Standards

- **Error Handling:** Never ignore errors. Propagate them with context or log
  them explicitly.
- **WASM:** Ensure targets are compatible with `js/wasm` where applicable.

## Workspace Conventions

- Use "Conventional Commits 1.0.0" for commit messages.
- Prefer rebase over merge: `git pull origin main --rebase`.
- PRs should be created against the `main` branch.

## CI/CD

- CI runs on GitHub Actions (Ubuntu).
- Uses Bazel caching (`~/.cache/bazel-disk-cache`,
  `~/.cache/bazel-repository-cache`).


## Specification

* Record specification to `//ai/spec.md`.

