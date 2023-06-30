FROM alpine:latest as cert-generator

RUN apk --no-cache add openssl

WORKDIR /certs

RUN openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 \
-subj "/C=DE/ST=Berlin/L=Berlin/O=TrackingDetector/CN=localhost" \
-keyout private.key -out certificate.crt

FROM nginx:latest

COPY --from=cert-generator /certs/private.key /etc/nginx/certs/private.key
COPY --from=cert-generator /certs/certificate.crt /etc/nginx/certs/certificate.crt

COPY ./infra/api-gateway/nginx.conf /etc/nginx/conf.d/default.conf

CMD [ "nginx", "-g", "daemon off;" ]