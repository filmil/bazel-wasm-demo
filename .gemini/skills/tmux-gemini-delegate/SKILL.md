---
name: tmux-gemini-delegate
description: >
  Delegate a new feature implementation to a fresh Gemini CLI instance in a new
  tmux window using a dedicated git worktree and branch. Use when you want to
  parallelize work or isolate a complex feature implementation.
---

# Tmux Gemini Delegate

This skill automates the process of spawning a new Gemini CLI agent to handle a
specific feature request in a new git worktree.

## Workflow

1.  **Identify the Topic**: A short, alphanumeric string (e.g., `prettify-ui`,
    `add-auth`) used for the worktree directory and branch name.
2.  **Describe the Feature**: A clear, concise instruction of what the new
    agent should implement.
3.  **Delegate**: Execute the bundled script to create a new tmux window and
    start the new agent.

## Script Usage

```bash
bash tmux-gemini-delegate/scripts/delegate.sh <topic> "<feature_description>"
```

The script performs the following actions:
- Validates the environment (must be inside tmux).
- Generates a random nonce for the branch name.
- Calculates the worktree path (sibling to the current directory).
- Creates a new tmux window and starts a new `gemini` instance.
- Instructs the new instance to:
    1. Create a new git worktree.
    2. Checkout a fresh branch from `origin/main`.
    3. Implement the feature.
    4. Create and send a Pull Request.

## Example

```bash
bash tmux-gemini-delegate/scripts/delegate.sh bootstrap-ui "Prettify the UI using Bootstrap 5"
```
