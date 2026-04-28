# bazel-wasm-demo

[![Test](https://github.com/filmil/bazel-wasm-demo/actions/workflows/test.yml/badge.svg)](https://github.com/filmil/bazel-wasm-demo/actions/workflows/test.yml)
[![Publish BCR](https://github.com/filmil/bazel-wasm-demo/actions/workflows/publish-bcr.yml/badge.svg)](https://github.com/filmil/bazel-wasm-demo/actions/workflows/publish-bcr.yml)
[![Publish](https://github.com/filmil/bazel-wasm-demo/actions/workflows/publish.yml/badge.svg)](https://github.com/filmil/bazel-wasm-demo/actions/workflows/publish.yml)
[![Tag and Release](https://github.com/filmil/bazel-wasm-demo/actions/workflows/tag-and-release.yml/badge.svg)](https://github.com/filmil/bazel-wasm-demo/actions/workflows/tag-and-release.yml)


An example Bazel-managed Go WASM web application using gRPC for API definitions.

## Overview

This project demonstrates how to build a Go WebAssembly client and a Go HTTP
server using Bazel. The client and server communicate via a JSON API that is
automatically generated from a gRPC `Greeter` service definition. The demo
allows the client to send a JSON API request with a name and receive "Hello,
NAME" in response.

## Running the Application

1. Start the server:
   ```bash
   bazel run //server
   ```
2. Open your browser and navigate to `http://localhost:8080`.
