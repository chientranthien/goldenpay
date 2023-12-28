# TODO not use root user
FROM chientt1993/golang:1.19
RUN ln -sf /bin/bash /bin/sh

WORKDIR /wrk
COPY . /wrk
RUN cd /wrk
ARG G_ENV=dev
RUN if [[ "$G_ENV" = "prod" ]] ; then make build/__service ; else echo "building for dev env, only copy from local bin" ; fi

RUN mkdir tmp; mv bin/__service tmp/__service ; rm bin/* -r; mv tmp/__service bin/__service

CMD G_CONFIG=internal/service/__service/config/config.yaml ./bin/__service/exc

