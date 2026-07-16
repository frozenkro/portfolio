#!/bin/bash

set -euo pipefail

if [ -n "${TMUX:-}" ]; then
  SN="$(tmux display-message -p '#S')"
  IX="$(tmux display-message -p '#I')"
else 
  SN="portfolio"
  tmux new-session -d -s "$SN"
  IX=0
fi

tmux rename-window -t "$SN:$IX" 'host'
tmux new-window -t "$SN" -n edit
tmux new-window -t "$SN" -n hermes

tmux split-window -t "$SN:hermes" -h
tmux split-window -t "$SN:host" -h

tmux select-window -t "$SN:edit.0"

tmux send-keys -t "$SN:edit.0" 'nvim' C-m
tmux send-keys -t "$SN:hermes.0" 'hermes --tui' C-m
tmux send-keys -t "$SN:hermes.1" 'hermes --tui' C-m


if [ -n "${TMUX:-}" ]; then
  go build -tags dev && ./portfolio
else
  tmux send-keys -t "$SN:host" 'go build -tags dev && ./portfolio'
fi
