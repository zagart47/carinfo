![Static Badge](https://img.shields.io/badge/%D1%81%D1%82%D0%B0%D1%82%D1%83%D1%81-%D0%B3%D0%BE%D1%82%D0%BE%D0%B2-blue)
![Static Badge](https://img.shields.io/badge/GO-1.22+-blue)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/zagart47/carinfo)
![GitHub last commit (by committer)](https://img.shields.io/github/last-commit/zagart47/carinfo)
![GitHub forks](https://img.shields.io/github/forks/zagart47/carinfo)

# Filmoteca
REST API для управления базой данных автомобилей

## Содержание
- [Технологии](#технологии)
- [Использование](#использование)
- [Разработка](#разработка)
- [Contributing](#contributing)
- [FAQ](#faq)
- [To do](#to-do)
- [Команда проекта](#команда-проекта)

## Технологии
- [Golang](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)

## Использование
Запустить сервер с БД PostgreSQL, указать настройки в .env
Запустить сервис:
```powershell
git clone https://github.com/zagart47/carinfo.git
cd carinfo/cmd/carinfo
go run main.go
```
Использовать следующие эндпоинты:
```http request
GET /cars
```
```http request
GET /car/{id}
```
```http request
POST /car/new
```
```http request
PUT /car/edit/{id}
```
```http request
DELETE /car/delete/{id}
```
В каталоге docs содержится документация [openapi.yaml](docs%2Fopenapi.yaml) с описанием всех эндпоинтов

## Разработка

### Требования
Для установки и запуска проекта необходимы golang, docker и прямые руки.

## Contributing
Если у вас есть предложения или идеи по дополнению проекта или вы нашли ошибку, то пишите мне в tg: @zagart47

## FAQ
### Зачем ты разработал этот проект?
Это тестовое задание.

## To do
- [x] Выставить REST-методы.
- [x] Возможность получения данных с фильтрацией по всем полям и пагинацией.
- [x] Удаление по идентификатору.
- [x] Изменение одного или нескольких полей по идентификатору.
- [x] Добавление новых авто в формате:
```json
{
  "regNums": [ "X111XX777",
               "Y222YY888" ]
}
```
- [x] Обогащение информации по авто с помощью стороннего API (в настройках [.env](internal%2Fconfig%2F.env) необходимо указать URL стороннего API).

## Команда проекта
- [Артур Загиров](https://t.me/zagart47) — Golang Developer

