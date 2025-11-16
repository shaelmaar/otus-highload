#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    create user social_network with password 'secret';
    create database social_network owner social_network;
    create extension if not exists pg_trgm;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname social_network <<-EOSQL
    create extension if not exists pg_trgm;
EOSQL
