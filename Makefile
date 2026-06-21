.PHONY: frontend-dev frontend-build mock-server typespec-build typespec-install frontend-install install dev \
	backend-run backend-build backend-test backend-gen docker-build test-e2e test-e2e-ui

frontend-install:
	cd frontend && npm install

typespec-install:
	cd typespec && npm install

install: frontend-install typespec-install

typespec-build:
	cd typespec && npm run build

frontend-build:
	cd frontend && npm run build

mock-server:
	cd frontend && ./scripts/mock-server.sh

frontend-dev:
	cd frontend && npm run dev

# ── Backend (Go) ─────────────────────────────────────────
backend-run:
	cd backend && make run

backend-build:
	cd backend && make build

backend-test:
	cd backend && make test

backend-gen:
	cd backend && make gen

docker-build:
	docker build -t booking-app .

# Поднимает Go-бэкенд (:3000) и фронтенд (:5173)
dev:
	@echo "Starting backend on :3000 and frontend on :5173..."
	@cd backend && go run ./cmd/server &
	@sleep 2
	@cd frontend && npm run dev

test-e2e:
	cd frontend && npx playwright test

test-e2e-ui:
	cd frontend && npx playwright test --ui
