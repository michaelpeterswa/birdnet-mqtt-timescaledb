CREATE USER grafana
WITH
    PASSWORD 'grafana';

GRANT USAGE ON SCHEMA sensors TO grafana;

GRANT
SELECT
    ON sensors.birdnet TO grafana;

-- ALTER ROLE grafana SET search_path to sensors,public;