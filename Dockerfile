FROM busybox
WORKDIR /mailbox
COPY out/build /mailbox/

ENV COMP=mailbox

EXPOSE 8080
ENTRYPOINT ["./mailbox-svc-linux"]