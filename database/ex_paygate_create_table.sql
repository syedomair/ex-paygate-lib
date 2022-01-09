

CREATE TABLE "merchant" (
  "id" bigserial PRIMARY KEY,
  "key" varchar NOT NULL,
  "name" varchar NOT NULL
);

CREATE TABLE "ledger" (
  "id" bigserial PRIMARY KEY,
  "merchant_id" bigint NOT NULL,
  "approve_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "action_type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "approve" (
  "id" bigserial PRIMARY KEY,
  "merchant_id" bigint NOT NULL,
  "cc_number" bigint NOT NULL,
  "cc_expiry" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "approve_key" varchar NULL,
  "amount_balance" bigint NOT NULL,
  "status" integer NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);


ALTER TABLE "ledger" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchant" ("id");
ALTER TABLE "ledger" ADD FOREIGN KEY ("approve_id") REFERENCES "approve" ("id");
ALTER TABLE "approve" ADD FOREIGN KEY ("merchant_id") REFERENCES "merchant" ("id");

INSERT INTO "merchant" ("key", "name") VALUES ('KEY1', 'XYZ Store');
INSERT INTO "merchant" ("key", "name") VALUES ('KEY2', 'ABC Store');
