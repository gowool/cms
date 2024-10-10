SET statement_timeout = 0;

--==============================================================================
--bun:split

CREATE TABLE "sessions" (
    "token" text PRIMARY KEY,
    "data" bytea NOT NULL,
    "expiry" timestamptz NOT NULL
);

--bun:split

CREATE INDEX "sessions_expiry_idx" ON "sessions" ("expiry");
