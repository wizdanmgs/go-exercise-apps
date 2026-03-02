FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY sponge_rest_postgres /app/sponge_rest_postgres
RUN chmod +x /app/sponge_rest_postgres

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./sponge_rest_postgres", "-c", "configs/sponge_rest_postgres.yml"]
# if you use the Configuration Center, sponge_rest_postgres.yml is changed to the Configuration Center configuration.
#CMD ["./sponge_rest_postgres", "-c", "configs/sponge_rest_postgres.yml", "-enable-cc"]
