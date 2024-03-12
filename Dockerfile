FROM golang

WORKDIR /chain

COPY . .

RUN go build -o ./bin/ChainReaction
EXPOSE 80
ENV PORT="80"

ENTRYPOINT ./bin/ChainReaction
