FROM node:12.18.3 as client
ARG PublicURL
ARG BackendURL
WORKDIR /client
COPY client .
ENV PUBLIC_URL ${PublicURL}
ENV REACT_APP_HOST_URL ${BackendURL}
RUN npm i
RUN rm ./node_modules/@susisu/mte-kernel/dist/mte-kernel.mjs && npm run build

FROM golang:1.15.0 as server
WORKDIR /server
COPY . .
RUN go mod download && go run tools/gen_rsa_key.go && go build -o app

FROM ubuntu:18.04
WORKDIR /app
COPY --from=client /client/build build
COPY --from=server /server/app .
COPY --from=server /server/public.pem .
COPY --from=server /server/private.pem .
RUN mkdir -p files/temp

CMD [ "./app" ]
