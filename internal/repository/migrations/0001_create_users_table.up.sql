CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz,
  "deleted_at" timestamptz,
  "first_name" varchar(50),
  "last_name" varchar(50),
  "email" varchar unique NOT NULL,
  "password" text NOT NULL,
  "role" varchar NOT NULL DEFAULT 'user',
  "email_verified" bool NOT NULL DEFAULT false
);

CREATE INDEX ON "users" ("email");