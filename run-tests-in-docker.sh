#!/bin/bash

# Запускаем тесты внутри Docker сети
echo "Запуск тестов внутри Docker сети..."

# Создаем временный контейнер с curl для тестирования
docker run --rm --network=cinemaabyss-network curlimages/curl:latest sh -c "
echo '=== Testing Monolith Service ==='
curl -s http://monolith:8080/health
echo ''
curl -s http://monolith:8080/api/movies | head -c 200
echo '...'

echo ''
echo '=== Testing Movies Service ==='
curl -s http://movies-service:8081/api/movies/health
echo ''
curl -s http://movies-service:8081/api/movies | head -c 200
echo '...'

echo ''
echo '=== Testing Events Service ==='
curl -s http://events-service:8082/api/events/health
echo ''

echo ''
echo '=== Testing Proxy Service ==='
curl -s http://proxy-service:8000/health
echo ''
curl -s http://proxy-service:8000/api/movies | head -c 200
echo '...'
"
