# TODO not use root user
FROM golang:1.19
RUN ln -sf /bin/bash /bin/sh

WORKDIR /wrk
COPY . /wrk
RUN cd /wrk
ARG G_ENV=dev
RUN if [[ "$G_ENV" = "prod" ]] ; then make build/event_handler ; else echo "building for dev env, only copy from local bin" ; fi

RUN mkdir tmp; mv bin/event_handler tmp/event_handler ; rm bin/* -r; mv tmp/event_handler bin/event_handler

CMD G_CONFIG=internal/service/event_handler/config/config.yaml ./bin/event_handler/exc

