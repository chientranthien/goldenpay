# TODO not use root user
FROM chientt1993/golang:1.19
RUN ln -sf /bin/bash /bin/sh

WORKDIR /wrk
COPY . /wrk
RUN cd /wrk
ARG G_ENV=dev
RUN if [[ "$G_ENV" = "prod" ]] ; then make build/wallet ; else echo "building for dev env, only copy from local bin" ; fi

RUN mkdir tmp; mv bin/wallet tmp/wallet ; rm bin/* -r; mv tmp/wallet bin/wallet

CMD G_CONFIG=internal/service/wallet/config/config.yaml ./bin/wallet/exc

