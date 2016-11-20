CREATE TABLE links (
  id BIGSERIAL PRIMARY KEY,

  from_page_id BIGINT NOT NULL,
  to_page_id BIGINT NOT NULL,

  created TIMESTAMP DEFAULT current_timestamp,
  updated TIMESTAMP DEFAULT current_timestamp
);
