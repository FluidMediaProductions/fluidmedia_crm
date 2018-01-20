FROM scratch

COPY static /static
COPY templates /templates
COPY migrations /migrations

COPY fluidmedia_crm /

CMD ["/fluidmedia_crm"]