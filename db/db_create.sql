DO
$body$
  BEGIN

    -- Create dblink Extension to connect to database
    CREATE EXTENSION IF NOT EXISTS dblink;

    PERFORM dblink_connect('root_session', 'dbname='|| current_database() || ' user='|| current_user ||' password=password');

     -- Create Admin User if not exists
    IF NOT EXISTS (
        SELECT *
        FROM   pg_catalog.pg_user
        WHERE  usename = 'mailbox_admin') THEN
      PERFORM dblink_exec('root_session', $$CREATE ROLE mailbox_admin
                WITH PASSWORD 'admin_password'
                LOGIN
                SUPERUSER
                CREATEROLE;$$);
    END IF;

      -- Create database if not exists
      IF NOT EXISTS (
        SELECT 1
        FROM pg_database
        WHERE datname = 'mailbox') THEN
          PERFORM dblink_exec('root_session', $$CREATE DATABASE mailbox
            OWNER mailbox_admin
            ENCODING 'UTF8'
            TEMPLATE template0
            LC_COLLATE 'C'
            LC_CTYPE 'C';$$);
      END IF;

    -- Create App User if not exists
      IF NOT EXISTS (
        SELECT *
        FROM   pg_catalog.pg_user
        WHERE  usename = 'mailbox_user') THEN
          PERFORM dblink_exec('root_session', $$CREATE ROLE mailbox_user WITH PASSWORD 'mailbox_password' LOGIN;$$);
      END IF;

    -- End the Root Session
    PERFORM dblink_disconnect('root_session');
END;
$body$
