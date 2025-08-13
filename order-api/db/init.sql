CREATE TABLE IF NOT EXISTS partners (
  id INT AUTO_INCREMENT PRIMARY KEY,
  partner_id VARCHAR(100) NOT NULL UNIQUE,
  name VARCHAR(200) NOT NULL,
  secret_key VARCHAR(200) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS conversions (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  transaction_id VARCHAR(100) NOT NULL UNIQUE,
  partner_id VARCHAR(100) NOT NULL,
  partner_name VARCHAR(200) NOT NULL,
  sale_amount DECIMAL(12,2) NOT NULL,
  metadata JSON NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (partner_id) REFERENCES partners(partner_id)
);

-- Seed a test partner
INSERT IGNORE INTO partners (partner_id, name, secret_key) VALUES ('partner-a', 'Partner A', 'secret_for_partner_a');
