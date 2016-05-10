DO
$body$
  BEGIN

    -- Create dblink Extension to connect to database
    CREATE EXTENSION IF NOT EXISTS dblink;

    PERFORM dblink_connect('seed_session', 'dbname='|| 'mailbox' || ' user='|| 'mailbox_admin' ||' password=admin_password');

      PERFORM dblink_exec('seed_session', 
        $$COPY users (emp_id, name, phone_no, email) FROM './employees.csv' DELIMITER ',' CSV HEADER;$$);

    PERFORM dblink_exec('seed_session', $$INSERT into dealers (name, icon)
        values ('Amazon', 'some-url');$$);

    PERFORM dblink_exec('seed_session', $$INSERT into dealers (name, icon)
        values ('Fedex', 'some-url');$$);

    PERFORM dblink_exec('seed_session', $$INSERT into dealers (name, icon)
        values ('Flipkart', 'some-url');$$);

    PERFORM dblink_exec('seed_session', $$INSERT into dealers (name, icon)
        values ('Bluedart', 'some-url');$$);

    PERFORM dblink_exec('seed_session', $$INSERT into dealers (name, icon)
        values ('Myntra', 'some-url');$$);

    PERFORM dblink_disconnect('seed_session');
END;
$body$

