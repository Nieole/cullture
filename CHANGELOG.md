2020-05-11
---
alter table projects add country varchar(100) null;
alter table projects add province varchar(100) null;
alter table projects add city varchar(200) null;
alter table projects add district varchar(200) null;
create index projects_country_province_city_district_index on projects (country, province, city, district);
2020-05-08
---
INSERT INTO culture.tags (id, name, code, is_delete, created_at, updated_at) VALUES ('e69fbcc4-18b1-472a-882c-c7fd47bea510', '献礼70周年', '献礼70周年', 0, '2020-05-08 11:26:36.000', '2020-05-08 11:26:38.000')