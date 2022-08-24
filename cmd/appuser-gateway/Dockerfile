FROM uhub.service.ucloud.cn/entropypool/alpine:3.16

WORKDIR /app

ARG token_access_secret=
ENV LOGIN_TOKEN_ACCESS_SECRET=$token_access_secret

COPY appuser-gateway /app

ENTRYPOINT ["./appuser-gateway"]

EXPOSE 50280 50281 50282

CMD ["run"]
