# Nomad Stunts: landing page + admin panel in one Go binary.
# Data (content.json, history, uploaded media) lives in the /data volume — back it up.

FROM golang:1.27-alpine AS build
WORKDIR /src
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /nomad .

FROM alpine:3.22
RUN adduser -D -H app && mkdir -p /data && chown app /data
COPY --from=build /nomad /app/nomad
COPY index.html about.html team.html projects.html gallery.html contacts.html favicon.ico /app/site/
COPY assets /app/site/assets
COPY upload /app/site/upload
ENV SITE_DIR=/app/site DATA_DIR=/data PORT=8080
USER app
VOLUME /data
EXPOSE 8080
CMD ["/app/nomad"]
