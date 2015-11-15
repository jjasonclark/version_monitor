FROM alpine:latest
MAINTAINER Jason Clark <jason@jjcconsultingllc.com>

ADD ca-bundle.crt /etc/ssl/certs/ca-certificates.crt
ADD build/* /app/
CMD ["/app/version_monitor"]
