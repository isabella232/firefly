BEGIN;
CREATE TABLE ffimethods (
  seq               SERIAL          PRIMARY KEY,
  id                UUID            NOT NULL,
  interface_id      UUID            NOT NULL,
  namespace         VARCHAR(64)     NOT NULL,
  name              VARCHAR(1024)   NOT NULL,
  pathname          VARCHAR(1024)   NOT NULL,
  description       TEXT            NOT NULL,
  params            BYTEA           NOT NULL,
  returns           BYTEA           NOT NULL
);

CREATE UNIQUE INDEX ffimethods_pathname ON ffimethods(interface_id,pathname);
COMMIT;