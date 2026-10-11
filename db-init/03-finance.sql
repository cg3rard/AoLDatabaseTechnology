CREATE TABLE utility_tariffs (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  utility_type  ENUM('electricity_token','water','ipl') NOT NULL,
  rate_per_unit DECIMAL(12,2)   NOT NULL,
  admin_fee     DECIMAL(12,2)   NOT NULL DEFAULT 0.00,
  updated_at    DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_utility_tariffs_type (utility_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE invoices (
  id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  invoice_number    VARCHAR(50)     NOT NULL,
  unit_id           BIGINT UNSIGNED NOT NULL,
  invoice_type      ENUM('ipl','water','electricity_token','laundry','card_replacement','permit_deposit') NOT NULL,
  total_amount      DECIMAL(12,2)   NOT NULL,
  payment_status    ENUM('unpaid','pending_verification','paid','cancelled') NOT NULL DEFAULT 'unpaid',
  payment_proof_url TEXT            NULL,
  due_date          DATE            NOT NULL,
  paid_at           DATETIME        NULL,
  created_at        DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_invoices_number (invoice_number),
  CONSTRAINT fk_invoices_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE utility_records (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id       BIGINT UNSIGNED NOT NULL,
  utility_type  ENUM('water','electricity') NOT NULL,
  initial_meter DECIMAL(10,2)   NOT NULL,
  final_meter   DECIMAL(10,2)   NOT NULL,
  usage_total   DECIMAL(10,2)   GENERATED ALWAYS AS (final_meter - initial_meter) STORED,
  reading_date  DATE            NOT NULL,
  invoice_id    BIGINT UNSIGNED NULL,
  inputted_by   BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_utility_records_unit    FOREIGN KEY (unit_id)     REFERENCES units (id)    ON DELETE CASCADE,
  CONSTRAINT fk_utility_records_invoice FOREIGN KEY (invoice_id)  REFERENCES invoices (id) ON DELETE SET NULL,
  CONSTRAINT fk_utility_records_user    FOREIGN KEY (inputted_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE token_purchases (
  id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id     BIGINT UNSIGNED NOT NULL,
  invoice_id  BIGINT UNSIGNED NOT NULL,
  nominal     DECIMAL(12,2)   NOT NULL,
  token_code  VARCHAR(30)     NULL,
  status      ENUM('pending','success','failed') NOT NULL DEFAULT 'pending',
  created_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_token_purchases_unit    FOREIGN KEY (unit_id)    REFERENCES units (id)    ON DELETE CASCADE,
  CONSTRAINT fk_token_purchases_invoice FOREIGN KEY (invoice_id) REFERENCES invoices (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;