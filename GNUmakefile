default: testacc

DD_VERSION ?= 2.58.4
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

# Apply modern Go idiom fixes (excludes the generated internal/ddclient package)
.PHONY: modernize modernize-check
modernize:
	go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest -fix -test $$(go list ./... | grep -v '/internal/ddclient')

# Report modernization opportunities without changing files
modernize-check:
	go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest -test $$(go list ./... | grep -v '/internal/ddclient')

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

# Regenerate internal/ddclient/client.gen.go from openapi-specs/$(DD_VERSION)/defect_dojo.json.
#
# Procedure (see internal/ddclient/generate.go for details):
#   1. Copy the versioned spec to the repo root as defect_dojo.json (gitignored).
#   2. Run the go.mod-pinned oapi-codegen (never a global binary) with
#      -generate types,client,skip-fmt. skip-fmt is required: the raw output
#      contains invalid Go that gofmt rejects, which steps 3-4 clean up.
#   3. Delete invalid bare-nil enum constants (` = <nil>` lines). String enums
#      with the quoted value "<nil>" are valid Go and are kept.
#   4. Remove const blocks and Valid() methods for enum types aliased to
#      time.Time / openapi_types.Date (Go cannot have constants of those types).
#   5. goimports -w.
# All post-processing steps are idempotent.

DDCLIENT_GEN := internal/ddclient/client.gen.go
OAPI_CODEGEN_VERSION = $(shell go list -m -f '{{.Version}}' github.com/oapi-codegen/oapi-codegen/v2)

# awk program for step 4: pass 1 collects enum type names aliased to
# time.Time / openapi_types.Date, pass 2 drops their const blocks and
# Valid() methods (harmless no-op if the blocks are already gone).
define DDCLIENT_STRIP_AWK
NR==FNR { if ($$0 ~ /^type [A-Za-z_][A-Za-z0-9_]*[ \t]+(time\.Time|openapi_types\.Date)[ \t]*$$/) bad[$$2]=1; next }
skip_const { if ($$0 ~ /^\)/) skip_const=0; next }
skip_func  { if ($$0 ~ /^\}/) skip_func=0; next }
/^\/\/ Defines values for / { t=$$5; sub(/\.$$/, "", t); if (t in bad) { skip_const=1; next } }
/^\/\/ Valid indicates whether the value is a known member of the / { t=$$(NF-1); if (t in bad) { skip_func=1; next } }
{ print }
endef
export DDCLIENT_STRIP_AWK

.PHONY: regen-client
regen-client:
	@test -f openapi-specs/$(DD_VERSION)/defect_dojo.json || \
		{ echo "error: openapi-specs/$(DD_VERSION)/defect_dojo.json not found"; exit 1; }
	cp -f openapi-specs/$(DD_VERSION)/defect_dojo.json ./defect_dojo.json
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		-generate types,client,skip-fmt -package ddclient -o $(DDCLIENT_GEN) defect_dojo.json
	sed -i '/ = <nil>$$/d' $(DDCLIENT_GEN)
	awk "$$DDCLIENT_STRIP_AWK" $(DDCLIENT_GEN) $(DDCLIENT_GEN) > $(DDCLIENT_GEN).tmp
	mv -f $(DDCLIENT_GEN).tmp $(DDCLIENT_GEN)
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w $(DDCLIENT_GEN); \
	else \
		go run golang.org/x/tools/cmd/goimports@latest -w $(DDCLIENT_GEN); \
	fi
	go build ./internal/ddclient/
	@echo "Regenerated $(DDCLIENT_GEN) from DefectDojo $(DD_VERSION) spec"

# Run version compatibility checks (spec collection only)
.PHONY: dd-compat
dd-compat:
	bash scripts/dd-version-compat.sh

# Run compatibility checks with acceptance tests
.PHONY: dd-compat-test
dd-compat-test:
	bash scripts/dd-version-compat.sh --test
