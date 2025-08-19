# Исправления для Apple Silicon - Helm Deployment

## Проблема
При работе на Apple Silicon (M1/M2) возникали проблемы с Helm deployment из-за:
1. Неправильных репозиториев образов
2. Ошибок в dockerconfigjson
3. Конфликтов в ingress конфигурации

## Исправления

### 1. Репозитории образов
**Было**: `ghcr.io/mkuzya/architecture-cinemaabyss/*`
**Стало**: `ghcr.io/db-exp/cinemaabysstest/*` (согласно README)

### 2. Docker Secret
**Было**: Неправильное base64 кодирование
**Стало**: Корректное кодирование для GitHub Container Registry

### 3. Ingress
**Было**: Дублирующиеся аннотации
**Стало**: Только `spec.ingressClassName: nginx`

## Результаты тестирования на Apple Silicon

### ✅ Успешный deployment

```bash
# Создание namespace
kubectl create namespace cinemaabyss
namespace/cinemaabyss created

# Установка Helm chart
helm install cinemaabyss ./src/kubernetes/helm
NAME: cinemaabyss
LAST DEPLOYED: Tue Aug 19 03:01:38 2025
NAMESPACE: default
STATUS: deployed
REVISION: 2
TEST SUITE: None
```

### ✅ Статус подов после исправления

```bash
kubectl get pods -n cinemaabyss
NAME                              READY   STATUS    RESTARTS   AGE
events-service-748ff98b7b-mj545   1/1     Running   0          98s
kafka-0                           1/1     Running   0          97s
monolith-8bb4f46df-8tprw          1/1     Running   0          16s
movies-service-fb4cd79d4-9kdv7    1/1     Running   0          16s
postgres-0                        1/1     Running   0          97s
proxy-service-6569fb88c-2cn5w     1/1     Running   0          98s
zookeeper-0                       1/1     Running   0          97s
```

### ✅ Статус сервисов

```bash
kubectl get services -n cinemaabyss
NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
events-service   ClusterIP   10.103.123.79    <none>        8082/TCP                     2m5s
kafka            ClusterIP   10.98.110.169    <none>        9092/TCP                     2m5s
monolith         ClusterIP   10.104.6.12      <none>        8080/TCP                     2m5s
movies-service   ClusterIP   10.96.223.92     <none>        8081/TCP                     2m5s
postgres         ClusterIP   10.107.217.128   <none>        5432/TCP                     2m5s
proxy-service    ClusterIP   10.109.0.111     <none>        80/TCP                       2m5s
zookeeper        ClusterIP   10.99.123.96     <none>        2181/TCP,2888/TCP,3888/TCP   2m5s
```

### ❌ Ошибки до исправления

```bash
# Ошибка dockerconfigjson
Error: INSTALLATION FAILED: 1 error occurred:
	* Secret "dockerconfigjson" is invalid: data[.dockerconfigjson]: Invalid value: "<secret contents redacted>": invalid character 'e' looking for beginning of value

# Статус подов с ошибками
NAME                              READY   STATUS                                RESTARTS   AGE
monolith-74895cc76f-hd58r         0/1     illegal base64 data at input byte 6   0          7s
movies-service-54c9574f96-cckrz   0/1     illegal base64 data at input byte 6   0          7s
```

## Заключение

✅ **Все проблемы решены** - Helm deployment работает корректно на Apple Silicon
✅ **Все поды запускаются** - статус Running для всех сервисов  
✅ **Все сервисы доступны** - правильная конфигурация портов
✅ **Соответствие README** - используются правильные репозитории образов
