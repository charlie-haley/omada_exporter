lint:
	golangci-lint run

generate-metrics-table:
	sh ./scripts/metric-markdown-table.sh
