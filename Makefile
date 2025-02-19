GO_CMD = go
WGO_CMD = wgo
TEMPL_CMD = templ
WEBAPP_PATH = cmd/webapp/main.go
OUTPUT_BINARY = c1

.PHONY: run-webapp run-templ run build-webapp

run-webapp: 
	@echo "Running webapp in live mode..."
	@export $$(grep -v '^#' .env | xargs) && $(WGO_CMD) run $(WEBAPP_PATH)

run-templ:
	@echo "Running templ in watch mode..."
	@$(TEMPL_CMD) generate --watch

run:
	@echo "Running both templ and webapp concurrently..."
	@(export $$(grep -v '^#' .env | xargs) && \
	  $(TEMPL_CMD) generate --watch & \
	  $(WGO_CMD) run $(WEBAPP_PATH))

build-webapp:
	@echo "Generating templates and building webapp..."
	@$(TEMPL_CMD) generate
	@$(GO_CMD) build -o $(OUTPUT_BINARY) $(WEBAPP_PATH)

