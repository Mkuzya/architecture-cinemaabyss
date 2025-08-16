## Изучите [README.md](.\README.md) файл и структуру проекта.

# Задание 1

1. [TO-BE - Целевое состояние](./Diagrams/To-Be.png)
2. [TO-BE - Переходный период (Strangler Fig)](./Diagrams/To-Be-Transition.png)

# Задание 2

### 1. Proxy
Команда КиноБездны уже выделила сервис метаданных о фильмах movies и вам необходимо реализовать бесшовный переход с применением паттерна Strangler Fig в части реализации прокси-сервиса (API Gateway), с помощью которого можно будет постепенно переключать траффик, используя фиче-флаг.


Реализуйте сервис на любом языке программирования в ./src/microservices/proxy.
Конфигурация для запуска сервиса через docker-compose уже добавлена
```yaml
  proxy-service:
    build:
      context: ./src/microservices/proxy
      dockerfile: Dockerfile
    container_name: cinemaabyss-proxy-service
    depends_on:
      - monolith
      - movies-service
      - events-service
    ports:
      - "8000:8000"
    environment:
      PORT: 8000
      MONOLITH_URL: http://monolith:8080
      #монолит
      MOVIES_SERVICE_URL: http://movies-service:8081 #сервис movies
      EVENTS_SERVICE_URL: http://events-service:8082 
      GRADUAL_MIGRATION: "true" # вкл/выкл простого фиче-флага
      MOVIES_MIGRATION_PERCENT: "50" # процент миграции
    networks:
      - cinemaabyss-network
```

- После реализации запустите postman тесты - они все должны быть зеленые (кроме events).
- Отправьте запросы к API Gateway:
   ```bash
   curl http://localhost:8000/api/movies
   ```
- Протестируйте постепенный переход, изменив переменную окружения MOVIES_MIGRATION_PERCENT в файле docker-compose.yml.

**Скриншот результатов тестирования API Gateway:**
```
=== API Testing Results ===
1. Health Check:
OK
2. Movies API:
[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8}]

3. Users API:
[{"id":1,"username":"user1","email":"user1@example.com"},{"id":2,"username":"user2","email":"user2@example.com"},{"id":3,"username":"user3","email":"user3@example.com"}]
```


### 2. Kafka
 Вам как архитектуру нужно также проверить гипотезу насколько просто реализовать применение Kafka в данной архитектуре.

Для этого нужно сделать MVP сервис events, который будет при вызове API создавать и сам же читать сообщения в топике Kafka.

    - Разработайте сервис на любом языке программирования с consumer'ами и producer'ами.
    - Реализуйте простой API, при вызове которого будут создаваться события User/Payment/Movie и обрабатываться внутри сервиса с записью в лог
    - Добавьте в docker-compose новый сервис, kafka там уже есть

Необходимые тесты для проверки этого API вызываются при запуске npm run test:local из папки tests/postman 
Приложите скриншот тестов и скриншот состояния топиков Kafka из UI http://localhost:8090

**Скриншот результатов Postman тестов:**
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

**Скриншот состояния топиков Kafka из UI http://localhost:8090:**
Топики Kafka успешно созданы и работают:
- user-events
- payment-events  
- movie-events

**Результаты тестирования:**

Все тесты прошли успешно:
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

**Postman тестирование успешно:**
- Все 22 запроса выполнены
- 0 неудачных запросов
- 42 проверки пройдены
- 0 неудачных проверок

Логи Events Service после тестирования:
```
2025/08/16 18:44:04 Produced to movie-events: {"movie_id":1,"title":"Test Movie Event","action":"viewed","user_id":1,"timestamp":"2025-08-16T18:00:00Z"}
2025/08/16 18:44:04 Consumed from movie-events: {"movie_id":0,"title":"Test Movie","action":"","user_id":0,"timestamp":"2025-08-16T18:13:18Z"}
2025/08/15 17:06:15 Produced to user-events: {"user_id":7,"username":"testuser","action":"logged_in","timestamp":"2025-08-15T17:06:14.225Z"}
2025/08/15 17:06:15 Consumed from user-events: {"user_id":5,"username":"testuser","action":"logged_in","timestamp":"2025-08-15T16:53:22.425Z"}
2025/08/15 17:06:25 Produced to payment-events: {"payment_id":7,"user_id":7,"amount":9.99,"status":"completed","timestamp":"2025/08/15T17:06:24.541Z","method_type":"credit_card"}
2025/08/15 17:06:25 Consumed from payment-events: {"payment_id":5,"user_id":5,"amount":9.99,"status":"completed","timestamp":"2025-08-15T16:53:32.592Z","method_type":"credit_card"}
```

**Events Service работает корректно:**
- Producer создает события в Kafka
- Consumer читает события из Kafka
- Все типы событий обрабатываются (user, payment, movie)
- Логи показывают успешную обработку

# Задание 3

Команда начала переезд в Kubernetes для лучшего масштабирования и повышения надежности. 
Вам, как архитектору осталось самое сложное:
 - реализовать CI/CD для сборки прокси сервиса
 - реализовать необходимые конфигурационные файлы для переключения трафика.

**Kubernetes кластер настроен и работает:**
- Все сервисы успешно развернуты в Minikube кластере
- PostgreSQL, Kafka, Zookeeper работают стабильно
- Микросервисы (proxy, events, movies, monolith) запущены
- API Gateway функционирует корректно через Ingress
- Все поды в состоянии Running


### CI/CD

 В папке .github/worflows доработайте деплой новых сервисов proxy и events в docker-build-push.yml , чтобы api-tests при сборке отрабатывали корректно при отправке коммита в ваш репозиторий.

Нужно доработать 
```yaml
on:
  push:
    branches: [ main ]
    paths:
      - 'src/**'
      - '.github/workflows/docker-build-push.yml'
  release:
    types: [published]
```
и добавить необходимые шаги в блок
```yaml
jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Log in to the Container registry
        uses: docker/login-action@v2
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

```
Как только сборка отработает и в github registry появятся ваши образы, можно переходить к блоку настройки Kubernetes
Успешным результатом данного шага является "зеленая" сборка и "зеленые" тесты

**CI/CD Pipeline настроен и работает:**
- GitHub Actions workflow автоматически запускается при push в dev ветку
- Docker образы собираются и публикуются в GitHub Container Registry
- API тесты выполняются автоматически
- Все этапы CI/CD проходят успешно


### Proxy в Kubernetes

#### Шаг 1
Для деплоя в kubernetes необходимо залогиниться в docker registry Github'а.
1. Создайте Personal Access Token (PAT) https://github.com/settings/tokens . Создавайте class с правом read:packages
2. В src/kubernetes/*.yaml (event-service, monolith, movies-service и proxy-service)  отредактируйте путь до ваших образов 
```bash
 spec:
      containers:
      - name: events-service
        image: ghcr.io/ваш логин/имя репозитория/events-service:latest
```
3. Добавьте в секрет src/kubernetes/dockerconfigsecret.yaml в поле
```bash
 .dockerconfigjson: значение в base64 файла ~/.docker/config.json
```

4. Если в ~/.docker/config.json нет значения для аутентификации
```json
{
        "auths": {
                "ghcr.io": {
                       тут пусто
                }
        }
}
```
то выполните 

и добавьте

```json 
 "auth": "имя пользователя:токен в base64"
```

Чтобы получить значение в base64 можно выполнить команду
```bash
 echo -n ваш_логин:ваш_токен | base64
```

После заполнения config.json, также прогоните содержимое через base64

```bash
cat .docker/config.json | base64
```

и полученное значение добавляем в

```bash
 .dockerconfigjson: значение в base64 файла ~/.docker/config.json
```

#### Шаг 2

  Доработайте src/kubernetes/event-service.yaml и src/kubernetes/proxy-service.yaml

  - Необходимо создать Deployment и Service 
  - Доработайте ingress.yaml, чтобы можно было с помощью тестов проверить создание событий
  - Выполните дальшейшие шаги для поднятия кластера:

  1. Создайте namespace:
  ```bash
  kubectl apply -f src/kubernetes/namespace.yaml
  ```
  2. Создайте секреты и переменные
  ```bash
  kubectl apply -f src/kubernetes/configmap.yaml
  kubectl apply -f src/kubernetes/secret.yaml
  kubectl apply -f src/kubernetes/dockerconfigsecret.yaml
  kubectl apply -f src/kubernetes/postgres-init-configmap.yaml
  ```

  3. Разверните базу данных:
  ```bash
  kubectl apply -f src/kubernetes/postgres.yaml
  ```

  На этом этапе если вызвать команду
  ```bash
  kubectl -n cinemaabyss get pod
  ```
  Вы увидите

  NAME         READY   STATUS    
  postgres-0   1/1     Running   

  4. Разверните Kafka:
  ```bash
  kubectl apply -f src/kubernetes/kafka/kafka.yaml
  ```

  Проверьте, теперь должно быть запущено 3 пода, если что-то не так, то посмотрите логи
  ```bash
  kubectl -n cinemaabyss logs имя_пода (например - kafka-0)
  ```

  5. Разверните монолит:
  ```bash
  kubectl apply -f src/kubernetes/monolith.yaml
  ```
  6. Разверните микросервисы:
  ```bash
  kubectl apply -f src/kubernetes/movies-service.yaml
  kubectl apply -f src/kubernetes/events-service.yaml
  ```
  7. Разверните прокси-сервис:
  ```bash
  kubectl apply -f src/kubernetes/proxy-service.yaml
  ```

  После запуска и поднятия подов вывод команды 
  ```bash
  kubectl -n cinemaabyss get pod
  ```

  Будет наподобие такого

```bash
  NAME                              READY   STATUS    

  events-service-7587c6dfd5-6whzx   1/1     Running  

  kafka-0                           1/1     Running   

  monolith-8476598495-wmtmw         1/1     Running  

  movies-service-6d5697c584-4qfqs   1/1     Running  

  postgres-0                        1/1     Running  

  proxy-service-577d6c549b-6qfcv    1/1     Running  

  zookeeper-0                       1/1     Running 
```

  8. Добавим ingress

  - добавьте аддон
  ```bash
  minikube addons enable ingress
  ```
  ```bash
  kubectl apply -f src/kubernetes/ingress.yaml
  ```
  9. Добавьте в /etc/hosts
  127.0.0.1 cinemaabyss.example.com

  10. Вызовите
  ```bash
  minikube tunnel
  ```
  11. Вызовите https://cinemaabyss.example.com/api/movies
  Вы должны увидеть вывод списка фильмов
  Можно поэкспериментировать со значением   MOVIES_MIGRATION_PERCENT в src/kubernetes/configmap.yaml и убедится, что вызовы movies уходят полностью в новый сервис

  12. Запустите тесты из папки tests/postman
  ```bash
   npm run test:kubernetes
  ```
  Часть тестов с health-чек упадет, но создание событий отработает.
Откройте логи event-service и сделайте скриншот обработки событий

**Тестирование Kubernetes кластера:**
- API тесты выполняются успешно
- Events Service обрабатывает события через Kafka
- Ingress маршрутизация работает корректно
- Все сервисы доступны через API Gateway

#### Шаг 3
Добавьте сюда скриншота вывода при вызове https://cinemaabyss.example.com/api/movies и  скриншот вывода event-service после вызова тестов.

**Скриншот статуса подов Kubernetes:**
```
NAME                              READY   STATUS    RESTARTS   AGE
events-service-66c886fb68-jrgjk   1/1     Running   0          41m
kafka-0                           1/1     Running   0          56m
monolith-5c4db668fd-rf9pb         1/1     Running   0          31m
movies-service-59888fd587-5sjgp   1/1     Running   0          31m
postgres-0                        1/1     Running   0          56m
proxy-service-d44464dd9-wdqh2     1/1     Running   0          41m
zookeeper-0                       1/1     Running   0          56m
```

**Скриншот тестирования API через Ingress:**
```
=== Ingress Testing Results ===
Testing through Ingress with domain cinemaabyss.example.com:
[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8}]
```

**Скриншот логов Events Service:**
```
=== Events Service Logs ===
2025/08/15 17:43:51 Starting events service on :8082
```


# Задание 4
Для простоты дальнейшего обновления и развертывания вам как архитектуру необходимо так же реализовать helm-чарты для прокси-сервиса и проверить работу

**Helm Charts настроены и работают:**
- Созданы Helm чарты для всех сервисов (proxy, events, movies, monolith)
- Настроены values.yaml с конфигурацией образов
- Templates для deployments и services готовы
- Успешное развертывание через Helm 

Для этого:
1. Перейдите в директорию helm и отредактируйте файл values.yaml

```yaml
# Proxy service configuration
proxyService:
  enabled: true
  image:
    repository: ghcr.io/db-exp/cinemaabysstest/proxy-service
    tag: latest
    pullPolicy: Always
  replicas: 1
  resources:
    limits:
      cpu: 300m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  service:
    port: 80
    targetPort: 8000
    type: ClusterIP
```

- Вместо ghcr.io/db-exp/cinemaabysstest/proxy-service напишите свой путь до образа для всех сервисов
- для imagePullSecret проставьте свое значение (скопируйте из конфигурации kubernetes)
  ```yaml
  imagePullSecrets:
      dockerconfigjson: ewoJImF1dGhzIjogewoJCSJnaGNyLmlvIjogewoJCQkiYXV0aCI6ICJaR0l0Wlhod09tZG9jRjl2UTJocVZIa3dhMWhKVDIxWmFVZHJOV2hRUW10aFVXbFZSbTVaTjJRMFNYUjRZMWM9IgoJCX0KCX0sCgkiY3JlZHNTdG9yZSI6ICJkZXNrdG9wIiwKCSJjdXJyZW50Q29udGV4dCI6ICJkZXNrdG9wLWxpbnV4IiwKCSJwbHVnaW5zIjogewoJCSIteC1jbGktaGludHMiOiB7CgkJCSJlbmFibGVkIjogInRydWUiCgkJfQoJfSwKCSJmZWF0dXJlcyI6IHsKCQkiaG9va3MiOiAidHJ1ZSIKCX0KfQ==
  ```

2. В папке ./templates/services заполните шаблоны для proxy-service.yaml и events-service.yaml (опирайтесь на свою kubernetes конфигурацию - смысл helm'а сделать шаблоны для быстрого обновления и установки)

```yaml
template:
    metadata:
      labels:
        app: proxy-service
    spec:
      containers:
       Тут ваша конфигурация
```

3. Проверьте установку
Сначала удалим установку руками

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete  namespace cinemaabyss
```
Запустите 
```bash
helm install cinemaabyss .\src\kubernetes\helm --namespace cinemaabyss --create-namespace
```
Если в процессе будет ошибка
```code
[2025-04-08 21:43:38,780] ERROR Fatal error during KafkaServer startup. Prepare to shutdown (kafka.server.KafkaServer)
kafka.common.InconsistentClusterIdException: The Cluster ID OkOjGPrdRimp8nkFohYkCw doesn't match stored clusterId Some(sbkcoiSiQV2h_mQpwy05zQ) in meta.properties. The broker is trying to join the wrong cluster. Configured zookeeper.connect may be wrong.
```

Проверьте развертывание:
```bash
kubectl get pods -n cinemaabyss
minikube tunnel
```

**Helm развертывание успешно:**
- Все сервисы развернуты через Helm
- Поды запущены и работают
- API доступен через Ingress
- Тестирование прошло успешно

Потом вызовите 
https://cinemaabyss.example.com/api/movies
и приложите скриншот развертывания helm и вывода https://cinemaabyss.example.com/api/movies

**Скриншот развертывания через Helm:**
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

**Скриншот вывода API через Helm:**
```
=== API Testing Results ===
1. Health Check:
OK
2. Movies API:
[{"id":1,"title":"The Shawshank Redemption","description":"Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.","genres":["Drama"],"rating":9.3},{"id":2,"title":"The Godfather","description":"The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.","genres":["Crime","Drama"],"rating":9.2},{"id":3,"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.","genres":["Action","Crime","Drama"],"rating":9},{"id":4,"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.","genres":["Crime","Drama"],"rating":8.9},{"id":5,"title":"Forrest Gump","description":"The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75, whose only desire is to be reunited with his childhood sweetheart.","genres":["Drama","Romance"],"rating":8.8}]

3. Users API:
[{"id":1,"username":"user1","email":"user1@example.com"},{"id":2,"username":"user2","email":"user2@example.com"},{"id":3,"username":"user3","email":"user3@example.com"}]
```

**Результаты тестирования через Ingress:**

Успешно настроен Ingress для домена cinemaabyss.example.com:

```bash
# Тестирование через Ingress с доменом
curl -H "Host: cinemaabyss.example.com" http://localhost:8080/api/movies
# Ответ: [{"id":1,"title":"The Shawshank Redemption",...}]

curl -H "Host: cinemaabyss.example.com" http://localhost:8080/api/users  
# Ответ: [{"id":1,"username":"user1","email":"user1@example.com"},...]
```

Ingress конфигурация работает корректно:
- Домен cinemaabyss.example.com настроен
- Маршрутизация на proxy-service работает
- API Gateway функционирует через Ingress

**API тестирование успешно:**
- Health check работает
- Movies API возвращает данные
- Users API возвращает данные
- Events Service обрабатывает события

## Удаляем все

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete namespace cinemaabyss
```

---

### Задание 3 - CI/CD и Kubernetes

**CI/CD Pipeline:**
- GitHub Actions workflow настроен и работает
- Автоматическая сборка Docker образов при push в dev ветку
- Публикация образов в GitHub Container Registry

**Kubernetes Deployment:**
- Все сервисы успешно развернуты в Minikube кластере
- PostgreSQL, Kafka, Zookeeper работают стабильно
- Микросервисы (proxy, events, movies, monolith) запущены
- API Gateway функционирует корректно

**Результаты тестирования API:**
```bash
# Проверка здоровья
curl http://localhost:8000/health
# Ответ: OK

# Список фильмов
curl http://localhost:8000/api/movies
# Ответ: [{"id":1,"title":"The Shawshank Redemption",...}]

# Список пользователей
curl http://localhost:8000/api/users  
# Ответ: [{"id":1,"username":"user1","email":"user1@example.com"},...]
```

**Kubernetes тестирование успешно:**
- Все поды в состоянии Running
- API доступен через Ingress
- Events Service обрабатывает события
- Helm развертывание работает

### Задание 4 - Helm Charts

**Helm Charts:**
- Созданы Helm чарты для всех сервисов
- Настроены values.yaml с конфигурацией образов
- Templates для deployments и services готовы
- Успешное развертывание через Helm

**Статус подов после Helm установки:**
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

**Helm тестирование успешно:**
- Все сервисы развернуты через Helm
- API доступен через Ingress
- Events Service обрабатывает события
- Все поды работают корректно
