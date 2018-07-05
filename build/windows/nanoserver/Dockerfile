ARG SOURCE_TAG=latest

FROM microsoft/nanoserver:${SOURCE_TAG}

USER ContainerAdministrator

COPY dist /

WORKDIR /

EXPOSE 9000

ENTRYPOINT ["/portainer.exe"]
