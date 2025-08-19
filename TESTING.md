# Тестирование 

## Задание 2: Реализация прокси-сервиса

### Тест 2.1: Запуск сервисов
**Команда:**
```bash
docker compose up -d
```

**Результат:**
```
[+] Running 13/13
 ✔ Container cinemaabyss-postgres        Healthy
 ✔ Container cinemaabyss-monolith        Healthy
 ✔ Container cinemaabyss-movies-service  Healthy
 ✔ Container cinemaabyss-events-service  Healthy
 ✔ Container cinemaabyss-proxy-service   Healthy
 ✔ Container cinemaabyss-test-runner     Started
```

### Тест 2.2: Проверка API Gateway
**Команда:**
```bash
curl http://localhost:8000/api/movies
```

**Результат:**
```json
[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8},{"id":6,"title":"Test Movie 341","description":"A test movie created by automated tests","genres":["Action","Drama"],"rating":4.5},{"id":7,"title":"Microservice Test Movie 55","description":"A test movie created by automated tests for the microservice","genres":["Sci-Fi","Thriller"],"rating":4.8}]
```

### Тест 2.3: Проверка Monolith Service
**Команда:**
```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/users
curl http://localhost:8080/api/movies
```

**Результат:**
```json
{"status": true}
[{"id":1,"username":"user1","email":"user1@example.com"},{"id":2,"username":"user2","email":"user2@example.com"},{"id":3,"username":"user3","email":"user3@example.com"}]
[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8}]
```

### Тест 2.4: Проверка Movies Service
**Команда:**
```bash
curl http://localhost:8081/api/movies/health
curl http://localhost:8081/api/movies
```

**Результат:**
```json
{"status": true}
[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8}]
```

### Тест 2.5: Проверка Events Service
**Команда:**
```bash
curl http://localhost:8082/api/events/health
curl -X POST http://localhost:8082/api/events/movie -H "Content-Type: application/json" -d '{"movie_id":1,"title":"Test Movie","action":"viewed","user_id":1}'
```

**Результат:**
```json
{"status": true}
{"status": "success", "message": "Movie event created and processed"}
```

### Тест 2.6: Проверка Proxy Service
**Команда:**
```bash
curl http://localhost:8000/health
curl http://localhost:8000/api/users
```

**Результат:**
```json
{"status": true}
[{"id":1,"username":"user1","email":"user1@example.com"},{"id":2,"username":"user2","email":"user2@example.com"},{"id":3,"username":"user3","email":"user3@example.com"}]
```

### Тест 2.7: Запуск Postman тестов
**Команда:**
```bash
cd tests/postman
npm install
node run-tests.js --environment docker --timeout 30000
```

**Результат:**
```
┌─────────────────────────┬───────────────────┬──────────────────┐
│                         │          executed │           failed │
├─────────────────────────┼───────────────────┼──────────────────┤
│              iterations │                 1 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│                requests │                22 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│            test-scripts │                22 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│      prerequest-scripts │                 0 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│              assertions │                42 │                0 │
├─────────────────────────┴───────────────────┴──────────────────┤
│ total run duration: 31.5s                                      │
├────────────────────────────────────────────────────────────────┤
│ total data received: 5.95kB (approx)                           │
├────────────────────────────────────────────────────────────────┤
│ average response time: 1413ms [min: 2ms, max: 10.6s, s.d.: 3.5s] │
└────────────────────────────────────────────────────────────────┘
```

### Тест 2.8: Проверка Kafka UI
**Команда:**
Открыть в браузере: http://localhost:8090

**Результат:**
- Топики созданы: `user-events`, `payment-events`, `movie-events`
- Сообщения обрабатываются корректно

### Тест 2.9: Проверка логов Events Service
**Команда:**
```bash
docker compose logs events-service
```

**Результат:**
```
2025/08/15 17:43:51 Produced to movie-events: {"movie_id":1,"title":"Test Movie Event","action":"viewed","user_id":1,"timestamp":"2025-08-15T17:43:51Z"}
2025/08/15 17:43:51 Consumed from movie-events: {"movie_id":1,"title":"Test Movie Event","action":"viewed","user_id":1,"timestamp":"2025-08-15T17:43:51Z"}
2025/08/15 17:06:15 Produced to user-events: {"user_id":7,"username":"testuser","action":"logged_in","timestamp":"2025-08-15T17:06:14.225Z"}
2025/08/15 17:06:15 Consumed from user-events: {"user_id":5,"username":"testuser","action":"logged_in","timestamp":"2025-08-15T16:53:22.425Z"}
2025/08/15 17:06:25 Produced to payment-events: {"payment_id":7,"user_id":7,"amount":9.99,"status":"completed","timestamp":"2025/08/15T17:06:24.541Z","method_type":"credit_card"}
2025/08/15 17:06:25 Consumed from payment-events: {"payment_id":5,"user_id":5,"amount":9.99,"status":"completed","timestamp":"2025/08/15T16:53:32.592Z","method_type":"credit_card"}
```

---

## Задание 3: Реализация CI/CD-пайплайнов

### Тест 3.1: Создание namespace
**Команда:**
```bash
kubectl create namespace cinemaabyss
```

**Результат:**
```
namespace/cinemaabyss created
```

### Тест 3.2: Применение секретов и конфигураций
**Команда:**
```bash
kubectl apply -f src/kubernetes/configmap.yaml
kubectl apply -f src/kubernetes/secret.yaml
kubectl apply -f src/kubernetes/dockerconfigsecret.yaml
kubectl apply -f src/kubernetes/postgres-init-configmap.yaml
```

**Результат:**
```
configmap/cinemaabyss-config created
secret/cinemaabyss-secret created
secret/dockerconfigjson created
configmap/postgres-init-config created
```

### Тест 3.3: Развертывание базы данных
**Команда:**
```bash
kubectl apply -f src/kubernetes/postgres.yaml
```

**Результат:**
```
statefulset.apps/postgres created
service/postgres created
```

### Тест 3.4: Развертывание Kafka
**Команда:**
```bash
kubectl apply -f src/kubernetes/kafka/kafka.yaml
```

**Результат:**
```
statefulset.apps/kafka created
service/kafka created
statefulset.apps/zookeeper created
service/zookeeper created
```

### Тест 3.5: Развертывание сервисов
**Команда:**
```bash
kubectl apply -f src/kubernetes/monolith.yaml
kubectl apply -f src/kubernetes/movies-service.yaml
kubectl apply -f src/kubernetes/events-service.yaml
kubectl apply -f src/kubernetes/proxy-service.yaml
```

**Результат:**
```
deployment.apps/monolith created
service/monolith created
deployment.apps/movies-service created
service/movies-service created
deployment.apps/events-service created
service/events-service created
deployment.apps/proxy-service created
service/proxy-service created
```

### Тест 3.6: Проверка статуса подов
**Команда:**
```bash
kubectl get pods -n cinemaabyss
```

**Результат:**
```
NAME                              READY   STATUS    RESTARTS   AGE
events-service-748ff98b7b-mj545   1/1     Running   0          98s
kafka-0                           1/1     Running   0          97s
monolith-8bb4f46df-8tprw          1/1     Running   0          16s
movies-service-fb4cd79d4-9kdv7    1/1     Running   0          16s
postgres-0                        1/1     Running   0          97s
proxy-service-6569fb88c-2cn5w     1/1     Running   0          98s
zookeeper-0                       1/1     Running   0          97s
```

### Тест 3.7: Настройка Ingress
**Команда:**
```bash
minikube addons enable ingress
kubectl apply -f src/kubernetes/ingress.yaml
```

**Результат:**
```
ingress.networking.k8s.io/cinemaabyss-ingress created
```

### Тест 3.8: Тестирование через Ingress
**Команда:**
```bash
# Добавить в /etc/hosts: 127.0.0.1 cinemaabyss.example.com
minikube tunnel
curl https://cinemaabyss.example.com/api/movies
```

**Результат:**
```json
[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8}]
```

---

## Задание 4: Реализация Helm-чартов

### Локальное тестирование Helm (Apple Silicon)

**Проблема:** На Apple Silicon (M1/M2) возникали проблемы с Helm deployment из-за неправильных репозиториев образов и ошибок в dockerconfigjson.

**Решение:** Использовать правильные репозитории образов согласно README:
- Было: `ghcr.io/mkuzya/architecture-cinemaabyss/*`
- Стало: `ghcr.io/db-exp/cinemaabysstest/*`

### Тест 4.1: Очистка предыдущего deployment
**Команда:**
```bash
kubectl delete all --all -n cinemaabyss
kubectl delete namespace cinemaabyss
```

**Результат:**
```
namespace "cinemaabyss" deleted
```

### Тест 4.2: Установка Helm chart
**Команда:**
```bash
helm install cinemaabyss ./src/kubernetes/helm --namespace cinemaabyss --create-namespace
```

**Результат:**
```
NAME: cinemaabyss
LAST DEPLOYED: Tue Aug 19 03:01:38 2025
NAMESPACE: cinemaabyss
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

### Тест 4.3: Проверка статуса Helm deployment
**Команда:**
```bash
kubectl get pods -n cinemaabyss
```

**Результат:**
```
NAME                              READY   STATUS    RESTARTS   AGE
events-service-66c886fb68-jrgjk   1/1     Running   0          10m
kafka-0                           1/1     Running   0          25m
monolith-5c4db668fd-rf9pb         1/1     Running   0          40s
movies-service-59888fd587-5sjgp   1/1     Running   0          35s
postgres-0                        1/1     Running   0          25m
proxy-service-d44464dd9-wdqh2     1/1     Running   0          10m
zookeeper-0                       1/1     Running   0          25m
```

### Тест 4.4: Тестирование API через Helm
**Команда:**
```bash
curl https://cinemaabyss.example.com/api/movies
```

**Результат:**
```json
[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8}]
```

### Тест 4.5: Запуск тестов из папки tests/postman
**Команда:**
```bash
npm run test:kubernetes
```

**Результат:**
```
┌─────────────────────────┬───────────────────┬──────────────────┐
│                         │          executed │           failed │
├─────────────────────────┼───────────────────┼──────────────────┤
│              iterations │                 1 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│                requests │                22 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│            test-scripts │                22 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│      prerequest-scripts │                 0 │                0 │
├─────────────────────────┼───────────────────┼──────────────────┤
│              assertions │                42 │                0 │
├─────────────────────────┴───────────────────┴──────────────────┤
│ total run duration: 30.7s                                      │
├────────────────────────────────────────────────────────────────┤
│ total data received: 29.91kB (approx)                          │
├────────────────────────────────────────────────────────────────┤
│ average response time: 1375ms [min: 1ms, max: 10s, s.d.: 3.4s] │
└────────────────────────────────────────────────────────────────┘
```

---

## CI/CD Тестирование

### Дополнительные тесты Helm
**Описание:** В GitHub Actions добавлены дополнительные тесты для проверки Helm deployment с корректными данными из задания.

**Команда:** Автоматически выполняется в CI/CD pipeline

**Результат:** Helm deployment тестируется автоматически

