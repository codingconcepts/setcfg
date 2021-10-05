FROM alpine
WORKDIR /
COPY setcfg /
ENTRYPOINT [ "./setcfg" ]