FROM scratch
COPY awssso /usr/bin/awssso
ENTRYPOINT [ "/usr/bin/awssso" ]