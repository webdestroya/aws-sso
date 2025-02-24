
.PHONY: test-release
test-release:
	@go tool goreleaser release --skip publish,sign,docker --clean --snapshot

.PHONY: test
test:
	@go test -v -timeout 30s -tags nodev,testmode  ./...


.PHONY: outdated
outdated:
	@go list -u -m -f '{{if not .Indirect}}{{if .Update}}{{.}}{{end}}{{end}}' all
