# Crear archivo Dockerfile
echo "FROM node:18-alpine

# Crear directorio de trabajo
WORKDIR /app

# Copiar package.json y package-lock.json
COPY package*.json ./

# Instalar dependencias
RUN npm install --production

# Copiar el resto de la app
COPY . .

# Exponer puerto
EXPOSE 3000

# Comando para correr la app
CMD [\"node\", \"index.js\"]" > Dockerfile