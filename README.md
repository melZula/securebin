# Securebin

- Тема Е2. Разработка веб-приложения для генерации страниц хранения конфиденциальной информации
  Разработать веб-приложение (frontend – Vue/React, backend – PHP/Go/Java), позволяющее:

1. Вводить текстовую информацию и сохранять в виде картинки
2. Сохранять информацию в виде страницы с картинкой с заданным временем хранения и предоставлять пользователю адрес и пароль для доступа
3. Генерировать запароленную страницу для отображения картинки с указанным временем жизни или количеством открытий (и отображением ранее сделанных обращений)
4. При обращении пользователя отмечать в cookie-файле время и путь обращения

## Tools

- Go v1.13.8
- PostgreSQL v12.3-1
- Nginx v1.18.0
- Vue.js v2.6.11
- Materialaze v1.0.0

### Go packages

- Router - [gorilla/mux](https://github.com/gorilla/mux)
- Handlers - [gorilla/handlers](https://github.com/gorilla/handlers)
- Logger - [logrus](https://github.com/sirupsen/logrus)
- _.toml_ parser - [BurntSushi/toml](https://github.com/BurntSushi/toml)
- Password - [sethvargo/go-password](https://github.com/sethvargo/go-password/password)
- Image - [go-freetype](https://github.com/golang/freetype)

## Usage

You need to run Nginx and setup reverse proxy. API path is `/api/...`<br/>
`$ make` build app<br/>
`$ ./apiserver` run

## Database

You can use [go-migrate](https://github.com/golang-migrate/migrate) to run migration:<br/>
`$ migrate -path migrations -database "postgres://localhost/securebin&sslmod=disable" up`
