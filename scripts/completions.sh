#!/bin/sh
set -e
rm -rf completions
mkdir completions
for sh in bash zsh fish; do
	go run -tags nodev main.go completion "$sh" >"completions/awssso.$sh"
done
