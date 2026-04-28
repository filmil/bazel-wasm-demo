# Project Specifications & AI Instructions History

This file documents the iterative requests and modifications made to initialize the `bazel-wasm-demo` project.

## Initialization & Setup
- **Goal:** Initialize a Bazel-managed Go WASM web application using gRPC for API definitions.
- **Actions Taken:**
  - Initialized `GEMINI.md` with project guidelines, setting Go as the primary language, Bazel as the build system (pinned version 9.0.1 via mandates), and defined standards for logging and error handling.
  - Sourced base code from `filmil/bazel-experiments/tree/main/wasm`.
  - Copied GitHub workflows (`.github`) and BCR configurations (`.bcr`) from `filmil/bazel-go-basic`.
  - Replaced the previous module name `github.com/grahamjenson/bazel-golang-wasm-proto` with `github.com/filmil/bazel-wasm-demo` globally.

## API Simplification & gRPC Customization
- **Request:** Remove all customization from the demo web app, retaining only a gRPC API that allows the client side to send a JSON API request to get a name and return "Hello, NAME".
- **Actions Taken:**
  - Simplified `protos/api/api.proto` to contain a single `Greeter` service with a `SayHello` RPC.
  - Rewrote the generated `protoc` plugin configurations (`tools/protoc-gen-client` and `tools/protoc-gen-server`) to output HTTP bindings for `GreeterHTTPClientInterface` and `GreeterHTTPServerInterface` using standard `text/template` and `protogen` plugins.
  - Cleaned up AWS/EC2 specific code in the server (`server/main.go`) and WASM frontend (`wasm/main.go`), aligning them with the new "Hello, NAME" gRPC logic.
  - Symlinked the generated `*.pb.go`, `*.client.go`, and `*.server.go` files back into the `//protos/api` directory from `bazel-bin` output to ease development.

## Build Definitions & Bzlmod Enhancements
- **Request:** Change repo references from `@io_bazel_rules_go` to `@rules_go` and update the protobuf bazel dependency to `34.1`.
- **Actions Taken:**
  - Replaced all usages of `@io_bazel_rules_go` with `@rules_go` inside `protos/api/BUILD.bazel`.
  - Updated `MODULE.bazel` to depend on `protobuf` version `34.1` and ensured `bzlmod` external dependency mappings were correct (`bazel mod tidy`).
  - Added a generated, non-empty placeholder icon at `//static/icon.png`.

## Workflows & Documentation
- **Request:** Setup GitHub build, test, and BCR publishing workflows, and add status icons and a summary to `README.md`.
- **Actions Taken:**
  - Configured GitHub Actions CI via `.github/workflows/test.yml` using `setup-bazel` for caching.
  - Updated `README.md` to include:
    - Build status badges for `test.yml` and `publish-bcr.yml`.
    - An overview summarizing the WASM/gRPC architecture.
    - Instructions on how to run the newly refactored application (`bazel run //server`).

## Go-app Application Integration
- **Request:** Register a serving API for a go-app application wasm.
- **Actions Taken:**
  - Integrated `github.com/maxence-charriere/go-app/v10/pkg/app` into `server/main.go`.
  - Registered `app.Handler` to serve the WASM application.
  - Configured `wasm/BUILD.bazel` to target `js/wasm` for the client binary.
  - Added a shared `ui` package to hold common UI components for both client and server.
  - Added an integration test `server/server_test.go` to verify resource delivery (HTML, WASM, Icon, Favicon).

## GitHub Workflows & BCR Configuration
- **Request:** Adapt GitHub workflows and BCR configurations to match the project structure.
- **Actions Taken:**
  - Updated `tag-and-release.yml` to build and package the `//server:server` binary for multiple architectures.
  - Adjusted `publish.yml` and `publish-bcr.yml` to point to `bazelbuild/bazel-central-registry` with the `filmil/bazel-central-registry` fork.
  - Verified `.bcr/metadata.template.json` contains correct homepage and repository links.

## Go-app v10 Upgrade
- **Request:** Update go-app version from v9 to v10.
- **Actions Taken:**
  - Replaced all imports and build references of `go-app/v9` with `go-app/v10`.
  - Adjusted `app.Route` calls to use the factory function signature introduced in v10.
