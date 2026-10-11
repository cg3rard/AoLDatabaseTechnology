CREATE TABLE units (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    unit_number VARCHAR(20) NOT NULL,
    tower_block VARCHAR(20) NOT NULL,
    `floor` INT NOT NULL,
    status ENUM('vacant', 'occupied', 'under_renovation') NOT NULL DEFAULT 'vacant',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE KEY uq_units_unit_number (unit_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE users(
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    unit_id BIGINT UNSIGNED NULL,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(190) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    role ENUM('admin', 'resident', 'laundry_staff', 'security', 'technician') NOT NULL,
    resident_type ENUM('owner', 'tenant') NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    UNIQUE KEY uq_users_email (email),
    CONSTRAINT fk_users_unit FOREIGN KEY (unit_id) REFERENCES units (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;