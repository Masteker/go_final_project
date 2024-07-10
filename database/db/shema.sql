CREATE TABLE "scheduler" (
	"id" INTEGER PRIMARY KEY AUTOINCREMENT,
	"date" DATE NOT NULL DEFAULT "",
	"title" VARCHAR(128) NOT NULL DEFAULT "",
	"comment" TEXT NOT NULL DEFAULT "",
	"repeat" TEXT NOT NULL DEFAULT ""
);

CREATE INDEX date_index ON scheduler (date);