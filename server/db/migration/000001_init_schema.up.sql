CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
)

