ARG BASEIMAGE=debian:12-slim

FROM ${BASEIMAGE}

ARG APP_NAME=cago

ENV TZ=Asia/Shanghai \
    APP_NAME=$APP_NAME

# 安装证书
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY $APP_NAME .

COPY ./data ./data

RUN ls -l && chmod +x $APP_NAME

CMD ["sh", "-c", "./$APP_NAME"]
