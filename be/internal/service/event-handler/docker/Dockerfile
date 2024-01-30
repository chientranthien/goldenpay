# TODO not use root user
FROM chientt1993/golang:1.19
RUN ln -sf /bin/bash /bin/sh

WORKDIR /wrk
COPY . /wrk
RUN cd /wrk
ARG G_ENV=dev
RUN if [[ "$G_ENV" = "prod" ]] ; then make build/event-handler ; else echo "building for dev env, only copy from local bin" ; fi

RUN mkdir tmp; mv bin/event-handler tmp/event-handler ; rm bin/* -r; mv tmp/event-handler bin/event-handler

CMD G_CONFIG=internal/service/event-handler/config/container_config.yaml ./bin/event-handler/exc

