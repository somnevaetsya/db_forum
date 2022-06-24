FROM golang:latest AS build

ADD . /app
WORKDIR /app
RUN GOAMD64=v3 go build -o api ./cmd/

FROM ubuntu:latest
COPY . .

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get -y update && apt-get install -y postgresql && rm -rf /var/lib/apt/lists/*

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "create user somnevaetsya with superuser password 'password';" &&\
    createdb -O somnevaetsya forumdb &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/api/ .

EXPOSE 5000
USER root
ENV PGPASSWORD password
CMD service postgresql start && psql -h localhost -d forumdb -U somnevaetsya -p 5432 -a -q -f ./db/db.sql && ./api