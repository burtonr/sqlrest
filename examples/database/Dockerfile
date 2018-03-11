FROM microsoft/mssql-server-linux

RUN mkdir /data
WORKDIR /data

COPY ./setup .

RUN chmod +x /data/insert-data.sh

# Default SQL Server port
EXPOSE 1433

ENV ACCEPT_EULA=Y
ENV SA_PASSWORD=someSecret42!

CMD /bin/bash ./entrypoint.sh