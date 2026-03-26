DROP TABLE IF EXISTS "statistics" CASCADE;

CREATE TABLE "statistics" (
    "statisticId" int4 NOT NULL,
    "tgUserId" varchar(255) NOT NULL,
    "exercise" varchar(255) NOT NULL,
    "count" float8 NOT NULL,
    "params" jsonb,
    "createdAt" timestamp with time zone NOT NULL DEFAULT NOW(),
    "statusId" int4 NOT NULL,
    PRIMARY KEY("statisticId")
);

DROP INDEX IF EXISTS "statistics_partitioned_statusId";
DROP INDEX IF EXISTS "statistics_partitioned_createdAt";
DROP INDEX IF EXISTS "statistics_partitioned_tgUserId";
DROP INDEX IF EXISTS "statistics_partitioned_exercise";

CREATE INDEX "statistics_partitioned_statusId"
    ON statistics ("statusId");

CREATE INDEX "statistics_partitioned_createdAt"
    ON statistics ("createdAt");

CREATE INDEX "statistics_partitioned_tgUserId"
    ON statistics ("tgUserId");

CREATE INDEX "statistics_partitioned_exercise"
    ON statistics ("exercise");

