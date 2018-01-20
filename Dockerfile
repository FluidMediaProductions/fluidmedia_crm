FROM scratch

COPY static /
COPY templates /
COPY migrations /

COPY fluidmedia_crm /

CMD ["/fluidmedia_crm"]