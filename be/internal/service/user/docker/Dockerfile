# TODO not use root user
FROM golang:1.19
RUN ln -sf /bin/bash /bin/sh

WORKDIR /wrk
COPY . /wrk
RUN cd /wrk
ARG G_ENV=dev
RUN if [[ "$G_ENV" = "prod" ]] ; then make build/user ; else echo "building for dev env, only copy from local bin" ; fi

CMD ./bin/user/exc

