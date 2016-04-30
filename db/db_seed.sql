DO
$body$
  BEGIN

    -- Create dblink Extension to connect to database
    CREATE EXTENSION IF NOT EXISTS dblink;

    PERFORM dblink_connect('seed_session', 'dbname='|| 'mailbox' || ' user='|| 'mailbox_admin' ||' password=admin_password');

      PERFORM dblink_exec('seed_session', 
        $$COPY users (emp_id, name, phone_no, email) FROM './employees.csv' DELIMITER ',' CSV HEADER;$$);

    PERFORM dblink_disconnect('seed_session');
END;
$body$

