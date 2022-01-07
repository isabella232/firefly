ALTER TABLE blobs ADD size BIGINT;

ALTER TABLE data ADD blob_name VARCHAR(1024);
ALTER TABLE data ADD blob_size BIGINT;

CREATE INDEX data_blob_name ON data(blob_name);
CREATE INDEX data_blob_size ON data(blob_size);