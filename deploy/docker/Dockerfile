ARG BASEIMAGE=debian:12-slim

FROM ${BASEIMAGE}

ARG APP_NAME=cago

ENV TZ=Asia/Shanghai \
    APP_NAME=$APP_NAME

# 国内镜像
RUN sed -i 's/deb.debian.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apt/sources.list.d/debian.sources

# 安装证书
RUN apt update && \
    apt install -y ca-certificates && \
    update-ca-certificates

WORKDIR /app

COPY $APP_NAME .

COPY ./data ./data

RUN ls -l && chmod +x $APP_NAME

CMD ["sh", "-c", "./$APP_NAME"]
