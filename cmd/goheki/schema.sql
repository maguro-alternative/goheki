/*
出典

type 2次元 3次元
*/
CREATE TABLE IF NOT EXISTS source (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    type TEXT NOT NULL
);
/*
人物の内容
*/
CREATE TABLE IF NOT EXISTS entry (
    id SERIAL NOT NULL PRIMARY KEY,
    source_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    image TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_id) REFERENCES source (id)
);
/*
タグの内容
*/
CREATE TABLE IF NOT EXISTS tag (
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);
/*
タグを付ける
*/
CREATE TABLE IF NOT EXISTS entry_tag (
    id SERIAL NOT NULL PRIMARY KEY,
    entry_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (id),
    FOREIGN KEY (tag_id) REFERENCES tag (id)
);
/*
好感
*/
CREATE TABLE IF NOT EXISTS heki_radar_chart (
    entry_id INTEGER PRIMARY KEY,
    ai INTEGER DEFAULT 0,
    nu INTEGER DEFAULT 0,
    FOREIGN KEY (entry_id) REFERENCES entry (id)
);
/*スリーサイズ 身長体重含む*/
CREATE TABLE IF NOT EXISTS bwh (
    entry_id INTEGER PRIMARY KEY,
    bust INTEGER,
    waist INTEGER,
    hip INTEGER,
    height INTEGER,
    weight INTEGER,
    FOREIGN KEY (entry_id) REFERENCES entry (id)
);
/*髪の長さの種類*/
CREATE TABLE IF NOT EXISTS hairlength_type (
    id SERIAL PRIMARY KEY,
    length TEXT NOT NULL
);
/*髪の長さ*/
CREATE TABLE IF NOT EXISTS hairlength (
    entry_id INTEGER PRIMARY KEY,
    hairlength_type_id INTEGER NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (id),
    FOREIGN KEY (hairlength_type_id) REFERENCES hairlength_type (id)
);
/*髪の色の種類*/
CREATE TABLE IF NOT EXISTS haircolor_type (
    id SERIAL PRIMARY KEY,
    color TEXT NOT NULL
);
/*髪色*/
CREATE TABLE IF NOT EXISTS haircolor (
    entry_id INTEGER PRIMARY KEY,
    color_id INTEGER NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (id),
    FOREIGN KEY (color_id) REFERENCES haircolor_type (id)
);
/*髪型の種類*/
CREATE TABLE IF NOT EXISTS hairstyle_type (
    id SERIAL PRIMARY KEY,
    style TEXT NOT NULL
);
/*髪型*/
CREATE TABLE IF NOT EXISTS hairstyle (
    entry_id INTEGER PRIMARY KEY,
    style_id INTEGER NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (id),
    FOREIGN KEY (style_id) REFERENCES hairstyle_type (id)
);
/*性格の種類*/
CREATE TABLE IF NOT EXISTS personality_type (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL
);
/*性格*/
CREATE TABLE IF NOT EXISTS personality (
    entry_id INTEGER PRIMARY KEY,
    type_id INTEGER NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (id),
    FOREIGN KEY (type_id) REFERENCES personality_type (id)
);
/*urlリンク*/
CREATE TABLE IF NOT EXISTS link (
    id SERIAL NOT NULL PRIMARY KEY,
    entry_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    url TEXT NOT NULL,
    nsfw BOOLEAN NOT NULL DEFAULT FALSE,
    darkness BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (entry_id) REFERENCES entry (id)
);
CREATE TABLE IF NOT EXISTS eyecolor_type (
    id SERIAL NOT NULL PRIMARY KEY,
    color TEXT NOT NULL
);
/*目の色*/
CREATE TABLE IF NOT EXISTS eyecolor (
    entry_id INTEGER PRIMARY KEY,
    color_id INTEGER NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (id),
    FOREIGN KEY (color_id) REFERENCES eyecolor_type (id)
);