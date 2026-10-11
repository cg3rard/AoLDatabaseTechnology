CREATE TABLE access_cards (
  id                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id             BIGINT UNSIGNED NOT NULL,
  rfid_uid            VARCHAR(64)     NOT NULL,
  card_label          VARCHAR(50)     NOT NULL,
  holder_name         VARCHAR(100)    NOT NULL,
  authorized_floor    INT             NOT NULL,
  has_facility_access BOOLEAN         NOT NULL DEFAULT TRUE,
  status              ENUM('active','blocked','lost','deactivated') NOT NULL DEFAULT 'active',
  issued_at           DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  blocked_at          DATETIME        NULL,
  block_reason        TEXT            NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_access_cards_rfid_uid (rfid_uid),
  CONSTRAINT fk_access_cards_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE card_requests (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id       BIGINT UNSIGNED NOT NULL,
  requested_by  BIGINT UNSIGNED NOT NULL,
  request_type  ENUM('new_additional','lost_replacement','damaged_replacement') NOT NULL,
  lost_card_id  BIGINT UNSIGNED NULL,
  holder_name   VARCHAR(100)    NOT NULL,
  reason        TEXT            NULL,
  fee_amount    DECIMAL(12,2)   NOT NULL DEFAULT 0.00,
  status        ENUM('pending','approved','ready_for_pickup','completed','rejected') NOT NULL DEFAULT 'pending',
  processed_by  BIGINT UNSIGNED NULL,
  created_at    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_card_requests_unit      FOREIGN KEY (unit_id)      REFERENCES units (id)        ON DELETE CASCADE,
  CONSTRAINT fk_card_requests_requester FOREIGN KEY (requested_by) REFERENCES users (id),
  CONSTRAINT fk_card_requests_lost_card FOREIGN KEY (lost_card_id) REFERENCES access_cards (id) ON DELETE SET NULL,
  CONSTRAINT fk_card_requests_processor FOREIGN KEY (processed_by) REFERENCES users (id)       ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;