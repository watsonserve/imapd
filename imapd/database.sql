
CREATE TYPE gender_t AS ENUM (
    Seen = 1,
    Answered = 2,
    Flagged = 4,
    Deleted = 8,
    Draft = 16,
    Recent = 32
);


CREATE TABLE IF NOT EXISTS "mail_boxes" (
    box_id    UUID        PRIMARY KEY,
    box_name  VARCHAR(64),
    user_id   UUID,
    validCode INT,
    nextUID   INT
);

CREATE TABLE IF NOT EXISTS "mail" (
    sn seq
    mail_id    INT,
    valid_code INT,
    box_name   VARCHAR(64) UNIQUE,
    user_id    UUID,
    Seen       BOOLEAN,
    Answered   BOOLEAN,
    Flagged    BOOLEAN,
    Deleted    BOOLEAN,
    Draft      BOOLEAN,
    Recent     BOOLEAN,
    mail       TEXT
)