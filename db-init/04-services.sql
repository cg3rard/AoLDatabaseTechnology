CREATE TABLE packages (
  id                   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id              BIGINT UNSIGNED NOT NULL,
  recipient_name       VARCHAR(100)    NOT NULL,
  courier_name         VARCHAR(50)     NOT NULL,
  receipt_number       VARCHAR(100)    NULL,
  package_photo_url    TEXT            NOT NULL,
  status               ENUM('waiting_pickup','collected') NOT NULL DEFAULT 'waiting_pickup',
  received_by_staff_id BIGINT UNSIGNED NOT NULL,
  collected_at         DATETIME        NULL,
  received_at          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_packages_unit  FOREIGN KEY (unit_id)              REFERENCES units (id) ON DELETE CASCADE,
  CONSTRAINT fk_packages_staff FOREIGN KEY (received_by_staff_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE laundry_orders (
  id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id      BIGINT UNSIGNED NOT NULL,
  staff_id     BIGINT UNSIGNED NULL,
  invoice_id   BIGINT UNSIGNED NULL,
  service_type ENUM('kiloan_cuci_setrika','kiloan_cuci_kering','satuan_dry_clean') NOT NULL,
  weight_kg    DECIMAL(6,2)    NULL,
  item_count   INT             NULL,
  total_cost   DECIMAL(12,2)   NOT NULL DEFAULT 0.00,
  status       ENUM('order_placed','received','washing','ironing','ready_for_pickup','completed') NOT NULL DEFAULT 'order_placed',
  notes        TEXT            NULL,
  created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  completed_at DATETIME        NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_laundry_orders_unit    FOREIGN KEY (unit_id)    REFERENCES units (id)    ON DELETE CASCADE,
  CONSTRAINT fk_laundry_orders_staff   FOREIGN KEY (staff_id)   REFERENCES users (id)    ON DELETE SET NULL,
  CONSTRAINT fk_laundry_orders_invoice FOREIGN KEY (invoice_id) REFERENCES invoices (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE facilities (
  id                    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name                  VARCHAR(100)    NOT NULL,
  max_capacity          INT             NOT NULL,
  open_time             TIME            NOT NULL,
  close_time            TIME            NOT NULL,
  slot_duration_minutes INT             NOT NULL,
  is_active             BOOLEAN         NOT NULL DEFAULT TRUE,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE facility_bookings (
  id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  facility_id  BIGINT UNSIGNED NOT NULL,
  user_id      BIGINT UNSIGNED NOT NULL,
  booking_date DATE            NOT NULL,
  start_time   TIME            NOT NULL,
  end_time     TIME            NOT NULL,
  status       ENUM('confirmed','cancelled','completed') NOT NULL DEFAULT 'confirmed',
  created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_bookings_slot (facility_id, booking_date, start_time),
  CONSTRAINT fk_bookings_facility FOREIGN KEY (facility_id) REFERENCES facilities (id) ON DELETE CASCADE,
  CONSTRAINT fk_bookings_user     FOREIGN KEY (user_id)     REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;