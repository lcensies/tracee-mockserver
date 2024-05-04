FROM golang:1.21.0-alpine AS builder

ARG USER=appuser
ARG UID=10001  
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" "${USER}" 

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify && mkdir -p /app_dist/data

COPY . . 
RUN GOOS=linux go build -ldflags="-w -s" -o /bin/mockserv

FROM scratch AS release

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/mockserv /bin/mockserv
COPY --from=builder --chown=appuser:appuser /app_dist /app

ENV PORT=8080
ENV GIN_MODE=release
EXPOSE $PORT

USER $USER:$USER

ENTRYPOINT ["/bin/mockserv"]
