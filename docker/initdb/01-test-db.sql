-- Runs once, on first initialization of an empty pgdata volume.
-- The table tests truncate the movies table, so they need their own database.
CREATE DATABASE greenlight_test;
