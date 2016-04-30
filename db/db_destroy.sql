DO
$body$
  BEGIN
    -- Create dblink Extension to connect to database
    CREATE EXTENSION IF NOT EXISTS dblink;

    -- Start the Root Session
    PERFORM dblink_connect('root_session', 'dbname='|| current_database() || ' user='|| current_user ||' password=password');

    PERFORM pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'mailbox';

    PERFORM dblink_exec('root_session', 'REVOKE ALL PRIVILEGES on database mailbox FROM mailbox_admin');
    PERFORM dblink_exec('root_session', 'DROP DATABASE IF EXISTS mailbox');
    PERFORM dblink_exec('root_session', 'DROP ROLE mailbox_user;');
    PERFORM dblink_exec('root_session', 'DROP ROLE mailbox_admin');

    -- End the Root Session
    PERFORM dblink_disconnect('root_session');

  END
$body$
