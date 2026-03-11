CREATE TABLE notification_deliveries (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  notification_id UUID NOT NULL,
  channel VARCHAR(50) NOT NULL,
  status VARCHAR(50) DEFAULT 'pending',
  retry_count INT DEFAULT 0,
  last_attempt_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_notification
    FOREIGN KEY(notification_id) 
    REFERENCES notifications(id)
    ON DELETE CASCADE
)