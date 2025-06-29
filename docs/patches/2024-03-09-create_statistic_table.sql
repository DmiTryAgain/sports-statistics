-- =============================================================================
-- Diagram Name: sportStatistics
-- Created on: 11/8/2024 2:26:28 PM
-- Diagram Version:
-- =============================================================================

DROP TABLE IF EXISTS "statistics" CASCADE;

CREATE TABLE "statistics" (
                              "statisticId" Serial NOT NULL,
                              "tgUserId" varchar(255) NOT NULL,
                              "exercise" varchar(255) NOT NULL,
                              "count" float8 NOT NULL,
                              "params" jsonb,
                              "createdAt" Timestamp with time zone NOT NULL Default now(),
                              "statusId" integer NOT NULL,
                              PRIMARY KEY("statisticId")
);

drop index if exists "statistics_partitioned_statusId";
drop index if exists  "statistics_partitioned_createdAt";
drop index if exists "statistics_partitioned_tgUserId";
drop index if exists "statistics_partitioned_exercise";

create index "statistics_partitioned_statusId"
    on statistics ("statusId");

create index "statistics_partitioned_createdAt"
    on statistics ("createdAt");

create index "statistics_partitioned_tgUserId"
    on statistics ("tgUserId");

create index "statistics_partitioned_exercise"
    on statistics ("exercise");

