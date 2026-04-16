FROM --platform=$BUILDPLATFORM alpine:latest AS certloader
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates

FROM scratch
ARG TARGETPLATFORM

# Copy CA Certificates
COPY --from=certloader /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY ${TARGETPLATFORM}/awssso /usr/bin/awssso
ENTRYPOINT [ "/usr/bin/awssso" ]