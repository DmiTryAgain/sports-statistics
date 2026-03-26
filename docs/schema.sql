-- =============================================================================
-- Diagram Name: sportStatistics
-- Created on: 3/26/2026 11:54:41 PM
-- Diagram Version: 
-- =============================================================================

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

CREATE INDEX "statistics_user_status_created" ON "statistics" (
	"tgUserId", 
	"statusId", 
	"createdAt"
);


CREATE INDEX "statistics_user_exercise_status" ON "statistics" (
	"tgUserId", 
	"exercise", 
	"statusId"
);




