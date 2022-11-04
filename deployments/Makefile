.PHONY: storage
storage:
	cd storage && sudo docker compose --project-name finances-bot up -d

.PHONY: logs
logs:
	mkdir -p logs/data
	touch logs/data/log.txt
	touch logs/data/offsets.yaml
	sudo chmod -R 777 ./logs/data
	cd logs && sudo docker compose --project-name finances-bot up -d

.PHONY: tracing
tracing:
	cd tracing && sudo docker compose --project-name finances-bot up -d

.PHONY: metrics
metrics:
	mkdir -p metrics/data
	sudo chmod -R 777 metrics/data
	cd metrics && sudo docker compose --project-name finances-bot up -d

all: logs tracing metrics