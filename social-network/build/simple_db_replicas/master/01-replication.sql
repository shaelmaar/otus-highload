create role replicator with replication login password 'repl_secret';

alter system set wal_level = 'replica';