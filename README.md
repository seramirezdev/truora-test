# Usando Go 1.12.7

## Project setup
Instalar las siguiente librerías al proyecto
```
go get github.com/go-chi/chi
go get github.com/lib/pq
go get github.com/likexian/whois-go
go get github.com/gocolly/colly
```

### Correr cockroachdb
En la raiz del proyecto ejecutar
```
cockroach start --insecure --http-addr localhost:8081
```

### Correr servidor GO
En la raiz del proyecto ejecutar
```
go run main.go
```

### Correr aplicación Vue-js
dentro de `/public`ejecutar
```
npm run serve
o
yarn serve
```
