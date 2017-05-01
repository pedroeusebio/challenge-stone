FROM library/postgres:9.6

ENV POSTGRES_USER docker
ENV POSTGRES_PASSWORD docker
ENV POSTGRES_DB docker

ADD ./config/postgres.sql /docker-entrypoint-initdb.d/