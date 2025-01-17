# TODO not use root user
FROM chientt1993/golang:1.22
RUN ln -sf /bin/bash /bin/sh

WORKDIR /wrk
COPY . /wrk
RUN cd /wrk
ARG G_ENV=dev
RUN if [[ "$G_ENV" = "prod" ]] ; then make build/http ; else echo "building for dev env, only copy from local bin" ; fi

RUN mkdir tmp; mv bin/http tmp/http ; rm bin/* -r; mv tmp/http bin/http

CMD G_CONFIG=internal/service/http/config/container_config.yaml ./bin/http/exc

