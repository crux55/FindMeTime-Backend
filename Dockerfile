FROM golang:1.18-alpine as fmtbe
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0
RUN apk add --no-cache libc6-compat 
RUN go build -o fmtbe ./cmd/FindMeTime/
EXPOSE 8080
ENTRYPOINT ["./fmtbe"]



FROM node:14 as migrations
WORKDIR /app
COPY package*.json ./
RUN npm install
RUN npm install db-migrate-pg
COPY . .
RUN git clone https://github.com/vishnubob/wait-for-it.git