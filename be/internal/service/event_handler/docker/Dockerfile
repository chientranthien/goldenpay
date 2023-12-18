# TODO not use root user
FROM golang:1.19

WORKDIR /wrk
COPY . /wrk
RUN cd /wrk
RUN make build/event_handler


CMD ./bin/event_handler/exc

