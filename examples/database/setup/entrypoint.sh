#!/bin/bash

set -e

run_cmd="./insert-data.sh"

>&2 echo "Allowing 30 seconds for SQL Server to bootstrap, then creating database.."
until $run_cmd & /opt/mssql/bin/sqlservr; do
>&2 echo "This should not be executing!"
done