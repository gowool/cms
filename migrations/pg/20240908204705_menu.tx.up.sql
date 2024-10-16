SET statement_timeout = 0;

--==============================================================================
--bun:split

CREATE TABLE "menus" (
    "id" integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    "node_id" integer REFERENCES "nodes"("id") ON DELETE SET NULL ,
    "name" varchar NOT NULL,
    "handle" varchar NOT NULL ,
    "enabled" boolean NOT NULL DEFAULT false,
    "created" timestamptz NOT NULL DEFAULT now(),
    "updated" timestamptz NOT NULL DEFAULT now()
);

--bun:split

CREATE INDEX "menus_created_updated_idx" ON "menus" ("created", "updated");
CREATE INDEX "menus_enabled_idx" ON "menus" ("enabled");

--bun:split

CREATE UNIQUE INDEX "menus_handle_unq" ON "menus" (lower("handle"));
