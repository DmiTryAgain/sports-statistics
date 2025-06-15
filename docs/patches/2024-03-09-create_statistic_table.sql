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



