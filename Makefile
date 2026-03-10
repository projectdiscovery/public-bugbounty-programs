DATA_IN := "src/data.yaml"
DATA_OUT := "dist/data.json"
DATA_SCHEMA := "src/data.schema.json"

lint:
	@yamllint $(DATA_IN)

compile: lint
	@yq eval -P -o=j $(DATA_IN) | jq -c > $(DATA_OUT)

test:
	@ajv test -s $(DATA_SCHEMA) -d $(DATA_OUT) --valid --all-errors --errors text

policy-checks:
	@TEMP_INVALID=$$(mktemp /tmp/invalid-XXXXXXXX.tmp); \
	go run ./cmd/httpx-tester/main.go -file "$(DATA_OUT)" -output "$$TEMP_INVALID" 2>/dev/null; \
	if [ -s "$$TEMP_INVALID" ]; then \
		echo "Invalid URLs found:"; \
		cat "$$TEMP_INVALID"; \
	fi; \
	rm -f "$$TEMP_INVALID"

validate-domains:
	@TEMP_INVALID=$$(mktemp /tmp/invalid-XXXXXXXX.tmp); \
	go run ./cmd/validate-domains/main.go -file "$(DATA_OUT)" -output "$$TEMP_INVALID" 2>/dev/null; \
	if [ -s "$$TEMP_INVALID" ]; then \
		cat "$$TEMP_INVALID"; \
	fi; \
	rm -f "$$TEMP_INVALID"

duplicate-domains:
	@TEMP_DUPLICATES=$$(mktemp /tmp/duplicates-XXXXXXXX.tmp); \
	jq -r '.programs[].domains[]' "$(DATA_OUT)" | sort | uniq -c | awk '$$1 > 1 { print $$2 }' > "$$TEMP_DUPLICATES"; \
	if [ -s "$$TEMP_DUPLICATES" ]; then \
		cat "$$TEMP_DUPLICATES"; \
		rm -f "$$TEMP_DUPLICATES"; \
		exit 1; \
	fi; \
	rm -f "$$TEMP_DUPLICATES"
