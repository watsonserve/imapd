
CREATE TYPE gender_t AS ENUM (
    Seen = 1,
    Answered = 2,
    Flagged = 4,
    Deleted = 8,
    Draft = 16,
    Recent = 32
);


CREATE TABLE IF NOT EXISTS "mail_boxes" (
    uidvalidity SERIAL      PRIMARY KEY,
    box_name    VARCHAR(64),
    user_id     UUID,
    recent      BIGINT,
    next_uid    INT
);

CREATE TABLE IF NOT EXISTS "mails" (
    mail_id     INT,
    uidvalidity INT,
    box_name    VARCHAR(64) UNIQUE,
    seen        BOOLEAN,
    answered    BOOLEAN,
    flagged     BOOLEAN,
    deleted     BOOLEAN,
    draft       BOOLEAN,
    mail        TEXT
)

/*
SELECT uidvalidity, recent FROM mail_boxes WHERE box_name=? AND user_id=?
SELECT COUNT(mail_id) AS exists FROM mails WHERE box_name=? AND uidvalidity=?
*/