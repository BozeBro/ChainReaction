FROM golang

WORKDIR /chain

COPY . .

RUN go build -o ./bin/ChainReaction
EXPOSE 3000
ENTRYPOINT ./bin/ChainReaction
