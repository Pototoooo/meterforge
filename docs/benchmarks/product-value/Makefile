COMPOSE := docker compose -p mf-product-value -f $(CURDIR)/compose.json
.PHONY: up test down
up:
	$(COMPOSE) up -d --wait --wait-timeout 180
test:
	python3 run.py
down:
	$(COMPOSE) down
