CREATE TABLE m_company (
    company_cd      INT(11)       DEFAULT 0 NOT NULL PRIMARY KEY,
    rr_cd           SMALLINT(6)   DEFAULT 0 NOT NULL,
    company_name    VARCHAR(256)  DEFAULT '' NOT NULL,
    company_name_k  VARCHAR(256)  DEFAULT NULL,
    company_name_h  VARCHAR(256)  DEFAULT NULL,
    company_name_r  VARCHAR(256)  DEFAULT NULL,
    company_url     VARCHAR(512)  DEFAULT NULL,
    company_type    SMALLINT(6)   DEFAULT 0,
    e_status        SMALLINT(6)   DEFAULT 0,
    e_sort          INT(11)       DEFAULT 0
);

CREATE INDEX m_company_rr_cd         ON m_company(rr_cd);
CREATE INDEX m_company_company_type  ON m_company(company_type);
CREATE INDEX m_company_e_sort        ON m_company(e_sort);

CREATE TABLE m_line (
    line_cd      INT(11)       DEFAULT 0 NOT NULL PRIMARY KEY,
    company_cd   INT(11)       DEFAULT 0 NOT NULL,
    line_name    VARCHAR(256)  DEFAULT '' NOT NULL,
    line_name_k  VARCHAR(256)  DEFAULT NULL,
    line_name_h  VARCHAR(256)  DEFAULT NULL,
    line_color_c VARCHAR(8)    DEFAULT NULL,
    line_color_t VARCHAR(32)   DEFAULT NULL,
    line_type    SMALLINT(6)   DEFAULT 0,
    lon          double        DEFAULT 0,
    lat          double        DEFAULT 0,
    zoom         SMALLINT(6)   DEFAULT 0,
    e_status     SMALLINT(6)   DEFAULT 0,
    e_sort       INT(11)       DEFAULT 0
);

CREATE INDEX m_line_company_cd ON m_line(company_cd);
CREATE INDEX m_line_e_sort     ON m_line(e_sort);

CREATE TABLE m_station (
    station_cd     INT(11)       DEFAULT 0 NOT NULL PRIMARY KEY,
    station_g_cd   INT(11)       DEFAULT 0 NOT NULL,
    station_name   VARCHAR(256)  DEFAULT '' NOT NULL,
    station_name_k VARCHAR(256)  DEFAULT NULL,
    station_name_r VARCHAR(256)  DEFAULT NULL,
    line_cd        INT(11)       DEFAULT 0 NOT NULL,
    pref_cd        SMALLINT(6)   DEFAULT 0,
    post           VARCHAR(32)   DEFAULT NULL,
    address        VARCHAR(1024) DEFAULT NULL,
    lon            DOUBLE        DEFAULT 0,
    lat            DOUBLE        DEFAULT 0,
    open_ymd       date          DEFAULT NULL,
    close_ymd      date          DEFAULT NULL,
    e_status       SMALLINT(6)   DEFAULT 0,
    e_sort         INT(11)       DEFAULT 0
);

CREATE INDEX m_station_station_g_cd ON m_station(station_g_cd);
CREATE INDEX m_station_line_cd      ON m_station(line_cd);
CREATE INDEX m_station_pref_cd      ON m_station(pref_cd);
CREATE INDEX m_station_e_sort       ON m_station(e_sort);

CREATE TABLE m_station_join (
    line_cd        INT(11)   DEFAULT 0 NOT NULL,
    station_cd1    INT(11)   DEFAULT 0 NOT NULL,
    station_cd2    INT(11)   DEFAULT 0 NOT NULL,
    PRIMARY KEY(line_cd,station_cd1,station_cd2)
);

