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
/*髪の長さ*/
CREATE TABLE IF NOT EXISTS hairlength (
    entry_id INTEGER PRIMARY KEY,
    length TEXT,
    FOREIGN KEY (entry_id) REFERENCES entry (id)
);
/*髪色*/
CREATE TABLE IF NOT EXISTS haircolor (
    entry_id INTEGER PRIMARY KEY,
    color TEXT,
    FOREIGN KEY (entry_id) REFERENCES entry (id)
);
/*髪型*/
CREATE TABLE IF NOT EXISTS hairstyle (
    entry_id INTEGER PRIMARY KEY,
    style TEXT,
    FOREIGN KEY (entry_id) REFERENCES entry (id)
);
/*性格*/
CREATE TABLE IF NOT EXISTS personality (
    id SERIAL NOT NULL PRIMARY KEY,
    entry_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY (entry_id) REFERENCES entry (id)
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