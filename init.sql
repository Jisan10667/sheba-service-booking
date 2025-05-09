-- filepath: /Users/jisanahmed/Desktop/sheba claude/.init.sql
CREATE DATABASE IF NOT EXISTS sheba_service_booking_db;
USE sheba_service_booking_db;

-- Users table
CREATE TABLE users (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  phone VARCHAR(20) DEFAULT NULL,
  role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  token VARCHAR(255) DEFAULT NULL,
  refresh_token VARCHAR(255) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY email (email)
);

-- Categories table
CREATE TABLE categories (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  parent_category_id INT DEFAULT NULL,
  is_active TINYINT(1) DEFAULT 1,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  icon_url VARCHAR(255) DEFAULT NULL,
  display_order INT DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY name (name),
  KEY parent_category_id (parent_category_id),
  CONSTRAINT categories_ibfk_1 FOREIGN KEY (parent_category_id) REFERENCES categories(id)
);

-- Services table
CREATE TABLE services (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10,2) NOT NULL,
  category_id INT NOT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  is_active TINYINT(1) DEFAULT NULL,
  is_featured TINYINT(1) DEFAULT 0,
  estimated_time_minutes INT DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_service_category (category_id),
  CONSTRAINT fk_service_category FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Bookings table
CREATE TABLE bookings (
  id INT NOT NULL AUTO_INCREMENT,
  service_id INT NOT NULL,
  user_name VARCHAR(255) NOT NULL,
  phone_number VARCHAR(20) NOT NULL,
  email VARCHAR(255) DEFAULT NULL,
  scheduled_at DATETIME DEFAULT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  user_id INT DEFAULT NULL,
  total_price DECIMAL(10,2) DEFAULT NULL,
  duration INT DEFAULT 1,
  booking_reference_code VARCHAR(50) DEFAULT NULL,
  status ENUM('pending','confirmed','in_progress','completed','cancelled') DEFAULT 'pending',
  notes TEXT,
  PRIMARY KEY (id),
  UNIQUE KEY booking_reference_code (booking_reference_code),
  KEY fk_booking_service (service_id),
  KEY idx_booking_user (user_id),
  KEY idx_booking_status (status),
  KEY idx_booking_reference (booking_reference_code),
  CONSTRAINT fk_booking_service FOREIGN KEY (service_id) REFERENCES services(id),
  CONSTRAINT fk_booking_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Booking status history table
CREATE TABLE booking_status_history (
  id INT NOT NULL AUTO_INCREMENT,
  booking_id INT NOT NULL,
  status ENUM('pending','confirmed','in_progress','completed','cancelled') NOT NULL,
  created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  is_active TINYINT(1) DEFAULT NULL,
  notes TEXT,
  created_by INT DEFAULT NULL,
  estimated_completion_time DATETIME DEFAULT NULL,
  PRIMARY KEY (id),
  KEY fk_status_history_booking (booking_id),
  KEY fk_status_history_user (created_by),
  CONSTRAINT fk_status_history_booking FOREIGN KEY (booking_id) REFERENCES bookings(id),
  CONSTRAINT fk_status_history_user FOREIGN KEY (created_by) REFERENCES users(id)
);