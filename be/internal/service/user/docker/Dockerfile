# TODO not use root user
FROM chientt1993/golang:1.19
RUN ln -sf /bin/bash /bin/sh

WORKDIR /wrk
COPY . /wrk
RUN cd /wrk
ARG G_ENV=dev
RUN if [[ "$G_ENV" = "prod" ]] ; then make build/user ; else echo "building for dev env, only copy from local bin" ; fi

RUN mkdir tmp; mv bin/user tmp/user ; rm bin/* -r; mv tmp/user bin/user

CMD G_CONFIG=internal/service/user/config/config.yaml ./bin/user/exc

