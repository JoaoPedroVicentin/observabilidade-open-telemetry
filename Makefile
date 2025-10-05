conf ?= .env
include $(conf)
export $(shell sed 's/=.*//' $(conf))



## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@rm -f coverage.*
	@docker image rm -f JoaoPedroVicentin/cep-api:v1
	@docker image rm -f JoaoPedroVicentin/weather-api:v1




## ----- COMPOSE
.PHONY: up
up: ## Put the compose-prd containers up
	@docker-compose up -d

.PHONY: down
down: ## Put the compose-prd containers down
	@docker-compose down



## ----- DEV
.PHONY: dev-up
dev-up: ## Put the compose-dev containers up
	@docker-compose -f docker-compose-dev.yaml up -d

.PHONY: dev-down
dev-down: ## Put the compose-dev containers down
	@docker-compose -f docker-compose-dev.yaml down

.PHONY: dev
dev: dev-down dev-up ## Restart compose-dev containers



## ---------- MAIN
.PHONY: run
run: ## Make requests to call-cep-api and weather-api
	@echo -e -----------------" cep-api -----------------"
	@echo -n "422: "; curl -s "http://localhost:8080/cep" -d '{"cep": "1234567"}'
	@echo -n "200: "; curl -s "http://localhost:8080/cep" -d '{"cep": "13330250"}'

	@echo -e "\n\n------------- weather-api -------------"
	@echo -n "422: "; curl -s "http://localhost:8081/cep/0100100"
	@echo -n "404: "; curl -s "http://localhost:8081/cep/01001009"
	@echo -n "200: "; curl -s "http://localhost:8081/cep/01001001"


.PHONY: call-cep-api
call-cep-api: ## Make a request to cep-api
	@echo -e -----------------" cep-api -----------------"
	@echo -n "200: "; curl -s "http://localhost:8080/cep" -d '{"cep": "13330250"}'

.PHONY: test
test: ## Run the tests
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html coverage.out -o coverage.html
