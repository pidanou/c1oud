GO_CMD = go
TEMPL_CMD = templ
WEBAPP_PATH = cmd/webapp/main.go
OUTPUT_BINARY = c1oud

.PHONY: run-webapp run-templ run build-webapp

run-webapp:
	@echo "Running webapp in live mode..."
	@$(GO_CMD) run $(WEBAPP_PATH) live

run-templ:
	@echo "Running templ in watch mode..."
	@$(TEMPL_CMD) generate --watch

run:
	@echo "Running both templ and webapp concurrently..."
	@$(TEMPL_CMD) generate --watch &
	@$(GO_CMD) run $(WEBAPP_PATH) live

build-webapp:
	@echo "Generating templates and building webapp..."
	@$(TEMPL_CMD) generate
	@$(GO_CMD) build -o $(OUTPUT_BINARY) $(WEBAPP_PATH)

