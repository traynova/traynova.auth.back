# Usa la imagen oficial de Golang como entorno de compilación (Builder)
FROM golang:1.25.0-alpine AS builder

# Configura variables de entorno para la construcción
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Crea el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de módulos de go
COPY go.mod go.sum ./

# Descarga todas las dependencias
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Compila la aplicación. Genera un ejecutable llamado "traynova-auth"
RUN go build -o traynova-auth main.go

# Crea una segunda etapa más ligera para el despliegue
FROM alpine:latest

# Instala certificados necesarios para llamadas HTTPS externas
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copia el binario compilado desde la etapa "builder"
COPY --from=builder /app/traynova-auth .

# Expone el puerto por defecto (asumido 8080 del stack web normal de backend)
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./traynova-auth"]