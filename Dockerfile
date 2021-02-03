FROM node:14 as client
ARG PublicURL
ARG BackendURL
ENV PUBLIC_URL ${PublicURL}
ENV REACT_APP_HOST_URL ${BackendURL}

COPY client/package.json client/package.json
COPY client/package-lock.json client/package-lock.json
RUN cd client && npm i && find node_modules -name '*.mjs' | xargs -i rm {}
COPY . .
RUN cd client && npm run build

FROM golang:1.15.0 as server
WORKDIR /server
COPY . .
RUN go mod download && go run tools/gen_rsa_key.go && go build -o app

FROM ubuntu:20.04
WORKDIR /app
COPY --from=client /client/build build
COPY --from=server /server/app .
COPY --from=server /server/public.pem .
COPY --from=server /server/private.pem .
RUN mkdir -p files/temp

CMD [ "./app" ]
