ARG BASEIMAGE=nginx:1.25

FROM $BASEIMAGE

COPY dist /usr/share/nginx/html

COPY deploy/docker/default.conf /etc/nginx/conf.d/default.conf
