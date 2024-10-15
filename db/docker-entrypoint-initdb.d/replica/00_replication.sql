CHANGE REPLICATION SOURCE TO
  SOURCE_HOST='mysql-primary',
  SOURCE_PORT=3306,
  SOURCE_USER='replica_user',
  SOURCE_PASSWORD='replica_password',
  SOURCE_SSL=1,
  SOURCE_AUTO_POSITION=1;
START REPLICA;


CREATE USER IF NOT EXISTS 'read_only_user'@'%' IDENTIFIED BY 'read_only_password';
GRANT SELECT ON *.* TO 'read_only_user'@'%';
