CREATE TABLE moving_renovation_permits (
  id                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id            BIGINT UNSIGNED NOT NULL,
  permit_type        ENUM('moving_in','moving_out','renovation') NOT NULL,
  start_date         DATE            NOT NULL,
  end_date           DATE            NOT NULL,
  contractor_name    VARCHAR(100)    NULL,
  worker_count       INT             NOT NULL DEFAULT 0,
  deposit_amount     DECIMAL(12,2)   NOT NULL DEFAULT 0.00,
  deposit_invoice_id BIGINT UNSIGNED NULL,
  status             ENUM('submitted','approved','rejected','in_progress','completed','deposit_refunded') NOT NULL DEFAULT 'submitted',
  notes              TEXT            NULL,
  created_at         DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_permits_unit    FOREIGN KEY (unit_id)            REFERENCES units (id)    ON DELETE CASCADE,
  CONSTRAINT fk_permits_invoice FOREIGN KEY (deposit_invoice_id) REFERENCES invoices (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE emergency_alerts (
  id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id      BIGINT UNSIGNED NOT NULL,
  triggered_by BIGINT UNSIGNED NOT NULL,
  alert_type   ENUM('medical','fire','security_threat') NOT NULL,
  status       ENUM('triggered','responding','resolved','false_alarm') NOT NULL DEFAULT 'triggered',
  handled_by   BIGINT UNSIGNED NULL,
  created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  resolved_at  DATETIME        NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_alerts_unit      FOREIGN KEY (unit_id)      REFERENCES units (id) ON DELETE CASCADE,
  CONSTRAINT fk_alerts_trigger   FOREIGN KEY (triggered_by) REFERENCES users (id),
  CONSTRAINT fk_alerts_handler   FOREIGN KEY (handled_by)   REFERENCES users (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE lost_and_found (
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  reported_by    BIGINT UNSIGNED NOT NULL,
  type           ENUM('lost','found') NOT NULL,
  item_name      VARCHAR(150)    NOT NULL,
  category       ENUM('electronics','wallet_card','keys','accessories','others') NOT NULL,
  location_found VARCHAR(150)    NOT NULL,
  photo_url      TEXT            NULL,
  status         ENUM('open','claimed','closed') NOT NULL DEFAULT 'open',
  claimed_by     BIGINT UNSIGNED NULL,
  claimed_at     DATETIME        NULL,
  created_at     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_lostfound_reporter FOREIGN KEY (reported_by) REFERENCES users (id),
  CONSTRAINT fk_lostfound_claimer  FOREIGN KEY (claimed_by)  REFERENCES users (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE announcements (
  id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  author_id    BIGINT UNSIGNED NOT NULL,
  title        VARCHAR(200)    NOT NULL,
  content      TEXT            NOT NULL,
  category     ENUM('maintenance','event','emergency','general') NOT NULL DEFAULT 'general',
  is_published BOOLEAN         NOT NULL DEFAULT TRUE,
  created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_announcements_author FOREIGN KEY (author_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;