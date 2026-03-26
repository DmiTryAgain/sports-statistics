DROP INDEX IF EXISTS "statistics_partitioned_statusId";
DROP INDEX IF EXISTS "statistics_partitioned_createdAt";
DROP INDEX IF EXISTS "statistics_partitioned_tgUserId";
DROP INDEX IF EXISTS "statistics_partitioned_exercise";

DROP INDEX IF EXISTS "statistics_user_status_created";
DROP INDEX IF EXISTS "statistics_user_exercise_status";
CREATE INDEX "statistics_user_status_created" ON "statistics" ("tgUserId", "statusId", "createdAt");
CREATE INDEX "statistics_user_exercise_status" ON "statistics" ("tgUserId", "exercise", "statusId");