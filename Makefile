lint:
	golangci-lint run

generate-metrics-table:
	go run main.go --host dummy --username dummy --password dummy mdocs > gen-metrics-table.md
