-- USERS
CREATE TABLE users
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- MICROGREENS LIBRARY
CREATE TABLE microgreens
(
    id                 INT AUTO_INCREMENT PRIMARY KEY,
    name               VARCHAR(255) NOT NULL,
    latin_name         VARCHAR(255),
    germination_days   INT,
    harvest_days       INT,
    optimal_temp       VARCHAR(50),
    light_requirements TEXT,
    humidity_level     VARCHAR(50),
    substrate          JSON,
    watering           TEXT,
    growth_notes       TEXT,
    tips               JSON,
    image_url          TEXT,
    is_popular         BOOLEAN DEFAULT FALSE
);

-- BATCHES (LOTS)
CREATE TABLE batches
(
    id                     INT AUTO_INCREMENT PRIMARY KEY,
    user_id                INT          NOT NULL,
    name                   VARCHAR(255) NOT NULL,
    microgreen_id          INT,
    sowing_date            DATE         NOT NULL,
    substrate              VARCHAR(255),
    comment                TEXT,
    estimated_harvest_days INT,
    created_at             DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (microgreen_id) REFERENCES microgreens (id)
);

-- OBSERVATIONS (PHENOLOGICAL JOURNAL)
CREATE TABLE observations
(
    id               INT AUTO_INCREMENT PRIMARY KEY,
    batch_id         INT  NOT NULL,
    date             DATE NOT NULL,
    note             TEXT,
    height_cm        FLOAT,
    water_status     ENUM ('none', 'light', 'normal', 'heavy') DEFAULT 'normal',
    light_type       ENUM ('natural', 'artificial', 'mixed')   DEFAULT 'natural',
    humidity_percent INT,
    created_at       DATETIME                                  DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (batch_id) REFERENCES batches (id)
);

-- OBSERVATION PHOTOS
CREATE TABLE observation_photos
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    observation_id INT  NOT NULL,
    photo_url      TEXT NOT NULL,
    label          ENUM ('default', 'before', 'after') DEFAULT 'default',
    created_at     DATETIME                            DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (observation_id) REFERENCES observations (id)
);


CREATE TABLE notify_tokens
(
    id      SERIAL PRIMARY KEY,
    user_id INT REFERENCES users (id),
    token   VARCHAR(255)
);

CREATE TABLE notify_history
(
    id       SERIAL PRIMARY KEY,
    user_id  INT REFERENCES users (id),
    title    VARCHAR(255),
    body     VARCHAR(255),
    sender   INT,
    receiver INT
);


-- ADVICE MESSAGES
CREATE TABLE advice_messages
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    microgreen_id INT,
    message       TEXT NOT NULL,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (microgreen_id) REFERENCES microgreens (id)
);

CREATE TABLE messages
(
    id          SERIAL PRIMARY KEY,
    sender_id   INTEGER   NOT NULL REFERENCES users (id),
    receiver_id INTEGER   NOT NULL REFERENCES users (id),
    text        TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reminders
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT  NOT NULL,
    message    TEXT NOT NULL,
    time       TIME NOT NULL,
    active     BOOLEAN  DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
