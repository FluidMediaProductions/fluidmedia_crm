FROM scratch

COPY static /static
COPY templates /templates
COPY migrations /migrations

COPY fluidmedia_crm /

EXPOSE 8080

CMD ["/fluidmedia_crm"]
