FROM amazonlinux

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY bin/movie-spots-api .

ENV PORT=3001

EXPOSE 3001

ENTRYPOINT ["/app/movie-spots-api"]
