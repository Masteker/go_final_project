CREATE TABLE "scheduler" (
	 PRIMARY KEY AUTOINCREMENT "id",
	date CHAR(8) NOT NULL DEFAULT " ",
	title VARCHAR(128) NOT NULL DEFAULT " ",
	comment	TEXT NOT NULL DEFAULT " ",
	repeat VARCHAR(128)NOT NULL DEFAULT " "
);

CREATE INDEX "scheduler_date" ON "scheduler" (
	"date"	DESC
);