BEGIN;
CREATE TABLE contractevents (
  seq              SERIAL          PRIMARY KEY,
  id               UUID            NOT NULL,
  namespace        VARCHAR(64)     NOT NULL,
  name             VARCHAR(1024)   NOT NULL,
  subscription_id  UUID            NOT NULL,
  outputs          TEXT,
  info             TEXT,
  timestamp        BIGINT          NOT NULL
);
CREATE INDEX contractevents_name ON contractevents(namespace,name);
CREATE INDEX contractevents_timestamp ON contractevents(timestamp);
CREATE INDEX contractevents_subscription_id ON contractevents(subscription_id);
COMMIT;