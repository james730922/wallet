-- Preflight: this query must return zero rows. Existing duplicates are financial
-- records and must be reconciled explicitly; this migration never deletes them.
SELECT record_id, COUNT(*) AS order_count
FROM order_scanpay
GROUP BY record_id
HAVING COUNT(*) > 1;

-- Apply only after the preflight result is empty.
ALTER TABLE order_scanpay
  DROP INDEX idx_order_scanpay_record,
  ADD UNIQUE KEY uk_order_scanpay_record (record_id);
