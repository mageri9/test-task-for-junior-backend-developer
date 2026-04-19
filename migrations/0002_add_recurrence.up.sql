ALTER TABLE tasks ADD COLUMN IF NOT EXISTS reccurence JSONB;

CREATE INDEX IF NOT EXISTS idx_tasks_reccurence ON tasks USING GIN (reccurence);

COMMENT ON COLUMN tasks.reccurence IS 'Настройки периодичности задачи (daily, monthly, specific_dates, parity)';