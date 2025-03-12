FROM scratch

COPY redfish_exporter /redfish_exporter

ENTRYPOINT ["/redfish_exporter"]
