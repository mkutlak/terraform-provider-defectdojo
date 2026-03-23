default: testacc

DD_VERSION ?= 2.54.3
export DD_VERSION

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m -parallel=4

.PHONY: generate-docs
generate-docs:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run ./...

# Start local DefectDojo for acceptance tests
.PHONY: dd-up dd-wait dd-down dd-logs

dd-up:
	docker compose up -d
	@$(MAKE) dd-wait

dd-wait:
	@echo "Waiting for DefectDojo API to be ready..."
	@for i in $$(seq 1 60); do \
		if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v2/ 2>/dev/null | grep -qE "200|401|403"; then \
			echo "DefectDojo is ready at http://localhost:8080"; \
			echo "Username: admin / Password: testpassword"; \
			exit 0; \
		fi; \
		echo "  Attempt $$i/60 - waiting 10s..."; \
		sleep 10; \
	done; \
	echo "DefectDojo failed to start"; \
	docker compose logs uwsgi nginx; \
	exit 1

dd-down:
	docker compose down -v

dd-logs:
	docker compose logs -f uwsgi

# Run acceptance tests against local Docker instance
.PHONY: testacc-local
testacc-local:
	DEFECTDOJO_BASEURL=http://localhost:8080 \
	DEFECTDOJO_USERNAME=admin \
	DEFECTDOJO_PASSWORD=testpassword \
	TF_ACC=1 go test ./internal/provider/ -v $(TESTARGS) -timeout 120m -parallel=4

# Fetch OpenAPI spec from running DefectDojo instance
.PHONY: dd-spec
dd-spec:
	@echo "Fetching OpenAPI spec from DefectDojo $(DD_VERSION)..."
	@mkdir -p openapi-specs/$(DD_VERSION)
	@TOKEN=$$(curl -sf -X POST http://localhost:8080/api/v2/api-token-auth/ \
		-H 'Content-Type: application/json' \
		-d '{"username":"admin","password":"testpassword"}' | \
		python3 -c "import sys,json; print(json.load(sys.stdin)['token'])") && \
	curl -sf http://localhost:8080/api/v2/oa3/schema/?format=json \
		-H "Authorization: Token $$TOKEN" \
		-o openapi-specs/$(DD_VERSION)/defect_dojo.json && \
	echo "Saved to openapi-specs/$(DD_VERSION)/defect_dojo.json"

# Run version compatibility checks (spec collection only)
.PHONY: dd-compat
dd-compat:
	bash scripts/dd-version-compat.sh

# Run compatibility checks with acceptance tests
.PHONY: dd-compat-test
dd-compat-test:
	bash scripts/dd-version-compat.sh --test
