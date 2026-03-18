default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

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
	TF_ACC=1 go test ./internal/provider/ -v $(TESTARGS) -timeout 120m
