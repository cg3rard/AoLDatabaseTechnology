CREATE TABLE maintenance_tickets (
  id                     BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id                BIGINT UNSIGNED NULL,
  reported_by             BIGINT UNSIGNED NOT NULL,
  assigned_technician_id BIGINT UNSIGNED NULL,
  category               ENUM('plumbing','electrical','air_conditioner','civil_structure') NOT NULL,
  title                  VARCHAR(150)    NOT NULL,
  description            TEXT            NOT NULL,
  photo_url              TEXT            NULL,
  priority               ENUM('low','medium','high','urgent') NOT NULL DEFAULT 'medium',
  status                 ENUM('open','assigned','in_progress','resolved','closed') NOT NULL DEFAULT 'open',
  resolution_notes       TEXT            NULL,
  created_at             DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  resolved_at            DATETIME        NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_tickets_unit       FOREIGN KEY (unit_id)                REFERENCES units (id) ON DELETE SET NULL,
  CONSTRAINT fk_tickets_reporter   FOREIGN KEY (reported_by)            REFERENCES users (id),
  CONSTRAINT fk_tickets_technician FOREIGN KEY (assigned_technician_id) REFERENCES users (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE visitor_passes (
  id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id        BIGINT UNSIGNED NOT NULL,
  visitor_name   VARCHAR(100)    NOT NULL,
  visitor_phone  VARCHAR(20)     NULL,
  vehicle_number VARCHAR(20)     NULL,
  pass_code      VARCHAR(32)     NOT NULL,
  valid_date     DATE            NOT NULL,
  check_in_time  DATETIME        NULL,
  check_out_time DATETIME        NULL,
  status         ENUM('active','checked_in','completed','expired') NOT NULL DEFAULT 'active',
  created_at     DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_visitor_passes_code (pass_code),
  CONSTRAINT fk_visitor_passes_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE parking_slots (
  id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  slot_code    VARCHAR(20)     NOT NULL,
  zone         VARCHAR(30)     NOT NULL,
  vehicle_type ENUM('car','motorcycle') NOT NULL,
  is_occupied  BOOLEAN         NOT NULL DEFAULT FALSE,
  PRIMARY KEY (id),
  UNIQUE KEY uq_parking_slots_code (slot_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE vehicles (
  id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  unit_id         BIGINT UNSIGNED NOT NULL,
  parking_slot_id BIGINT UNSIGNED NULL,
  plate_number    VARCHAR(20)     NOT NULL,
  vehicle_type    ENUM('car','motorcycle') NOT NULL,
  brand_model     VARCHAR(100)    NULL,
  status          ENUM('active','inactive') NOT NULL DEFAULT 'active',
  created_at      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_vehicles_plate (plate_number),
  CONSTRAINT fk_vehicles_unit FOREIGN KEY (unit_id)         REFERENCES units (id)         ON DELETE CASCADE,
  CONSTRAINT fk_vehicles_slot FOREIGN KEY (parking_slot_id) REFERENCES parking_slots (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;