CREATE TABLE IF NOT EXISTS movements (
    ID bigserial PRIMARY KEY,
	CreatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	Name text NOT NULL,
	Description text NOT NULL,
	Image text NOT NULL,
	Tutorials text[] NOt NULL,
	Skilltype text[] NOt NULL,
	Muscles text[] NOt NULL,
	Difficulty text NOt NULL,
	Equipments text[] NOt NULL,
	Prerequisite text[],
	version integer NOT NULL DEFAULT 1
);