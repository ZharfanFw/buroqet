version: "3.8"
services:
  postgres-test:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: testuser
      POSTGRES_PASSWORD: testpass
      POSTGRES_DB: settlement_test
    ports:
      - "5434:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U testuser -d settlement_test"]
      interval: 5s
      timeout: 5s
      retries: 5
      