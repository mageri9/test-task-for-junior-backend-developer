
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS source_task_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS scheduled_date DATE;


CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_source_scheduled
    ON tasks(source_task_id, scheduled_date)
    WHERE source_task_id IS NOT NULL AND scheduled_date IS NOT NULL;


COMMENT ON COLUMN tasks.source_task_id IS 'ID исходной задачи-шаблона (для периодических задач)';
COMMENT ON COLUMN tasks.scheduled_date IS 'Дата, на которую создана задача (для периодических)';