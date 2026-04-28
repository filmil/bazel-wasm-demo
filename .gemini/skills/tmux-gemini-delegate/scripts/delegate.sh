#!/bin/bash

# Delegate a feature request to a new Gemini instance in a tmux window.

TOPIC=$1
DESCRIPTION=$2

if [ -z "$TOPIC" ] || [ -z "$DESCRIPTION" ]; then
  echo "Usage: $0 <topic> \"<description>\""
  echo "Example: $0 prettify-ui \"Use bootstrap to make it look nice\""
  exit 1
fi

if [ -z "$TMUX" ]; then
  echo "Error: This script must be run inside a tmux session."
  exit 1
fi

# Generate a 6-character nonce
NONCE=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 6 | head -n 1)
BRANCH_NAME="ai-dev-$(date +%Y%m%d)-$NONCE-$TOPIC"
WORKTREE_PATH="../$TOPIC"

# The command to send to the new Gemini instance
COMMAND="gemini \"in a new git worktree at $WORKTREE_PATH, checkout a new branch $BRANCH_NAME from origin/main. Implement a new feature: $DESCRIPTION. When done, create and send a PR.\""

# Create a new window and get the new window ID
WINDOW_ID=$(tmux new-window -d -P -F "#{window_id}")

if [ -z "$WINDOW_ID" ]; then
  echo "Error: Failed to create new tmux window."
  exit 1
fi

# Send the command to the new window
tmux send-keys -t "$WINDOW_ID" "$COMMAND" Enter

echo "✅ Successfully delegated feature '$TOPIC' to new tmux window $WINDOW_ID."
echo "   Branch: $BRANCH_NAME"
echo "   Worktree: $WORKTREE_PATH"
