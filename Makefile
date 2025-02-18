
.PHONY: test-release
test-release:
	@go tool goreleaser release --skip publish --clean --snapshot