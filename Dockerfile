FROM node:14 as client
WORKDIR client

COPY ["client/package.json", "client/package-lock.json", "./"]
RUN npm i && find node_modules -name '*.mjs' | xargs -i rm {}
COPY client ./

ARG PublicURL
ARG BackendURL
ENV PUBLIC_URL ${PublicURL}
ENV REACT_APP_HOST_URL ${BackendURL}

RUN npm run build

FROM golang:1.15 as server
WORKDIR /server
COPY [ "go.mod", "go.sum", "./" ]
RUN go mod download
COPY . .
RUN go run tools/gen_rsa_key.go && go build -o app

FROM ubuntu:20.04
WORKDIR /app
RUN mkdir -p files/temp && apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=client /client/build build
COPY --from=server [ "/server/app", "/server/public.pem", "/server/private.pem", "./" ]

CMD [ "./app" ]
