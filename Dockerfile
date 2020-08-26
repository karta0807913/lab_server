FROM node:12.18.3 as client
WORKDIR /client
COPY client .
RUN npm i && npm build

FROM golang:1.15.0 as server
WORKDIR /app
COPY . .
RUN go run tools/gen_rsa_key.go
RUN go build -o app

FROM ubuntu:18.04

COPY --from=client /client/build build
COPY --from=server /server/app .
COPY --from=server /server/public.pem .
COPY --from=server /server/private.pem .

CMD [ "app" ]