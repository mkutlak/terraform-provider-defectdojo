default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: lint
lint:
	golangci-lint run ./...

# Start local DefectDojo for acceptance tests
.PHONY: dd-up dd-down dd-logs

dd-up:
	docker compose up -d
	@echo "Waiting for DefectDojo to initialize..."
	@sleep 60
	@echo "DefectDojo should be available at http://localhost:8080"
	@echo "Username: admin / Password: testpassword"

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
	TF_ACC=1 go test ./internal/provider/ -v $(TESTARGS) -timeout 120m
