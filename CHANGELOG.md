##2020-05-20
```sql
create table banners
(
    id         char(36)     not null
        primary key,
    content    text         null,
    title      varchar(500) null,
    sort       int          not null,
    is_delete  tinyint(1)   not null,
    target     varchar(200) null,
    source     varchar(200) null,
    created_at datetime     not null,
    updated_at datetime     not null
);

INSERT INTO `banners` VALUES ('1704ce99-36f5-45e4-b3c7-764a597a21e2', NULL, NULL, 2, 0, 'http://maka.im/danyeviewer/6707418/OVJR8VBDW6707418?hash=f3af79a657ce9c5061ec6e87c34ffa13&from=timeline', 'http://ggjc-oss.oss-cn-chengdu.aliyuncs.com/2020/05/20/16/04/29/mY8tHp2R5dCaTCCks7nZNxshdXP6cEEK.jpeg', '2020-05-20 10:04:10', '2020-05-20 10:04:11');
INSERT INTO `banners` VALUES ('28aec275-b470-4bbb-8311-f5e96c4152d3', NULL, NULL, 3, 0, 'http://maka.im/danyeviewer/6707418/OVJR8VBDW6707418?hash=f3af79a657ce9c5061ec6e87c34ffa13&from=timeline', 'http://ggjc-oss.oss-cn-chengdu.aliyuncs.com/2020/05/20/16/04/38/Y6RiRciYWkB8N3PXQZ7r8CmycnXM5Ge8.jpeg', '2020-05-20 10:04:24', '2020-05-20 10:04:26');
INSERT INTO `banners` VALUES ('726a32d4-4fd1-4781-bbed-e94fdd5e9d1c', NULL, NULL, 4, 0, 'http://maka.im/danyeviewer/6707418/OVJR8VBDW6707418?hash=f3af79a657ce9c5061ec6e87c34ffa13&from=timeline', 'http://ggjc-oss.oss-cn-chengdu.aliyuncs.com/2020/05/20/16/04/48/E7f2RCsSQG2QjBPFD4GX8YGKkQA4H7yR.jpeg', '2020-05-20 10:04:39', '2020-05-20 10:04:41');
INSERT INTO `banners` VALUES ('8c5d1bef-6952-4628-9596-5fff66042b3a', NULL, NULL, 1, 0, 'http://maka.im/danyeviewer/6707418/OVJR8VBDW6707418?hash=f3af79a657ce9c5061ec6e87c34ffa13&from=timeline', 'http://ggjc-oss.oss-cn-chengdu.aliyuncs.com/2020/05/20/16/04/11/m3ZFEFj3mANZzycdirS3Xc2ZKc54QTBe.jpeg', '2020-05-20 10:03:52', '2020-05-20 10:03:52');
```
##2020-05-11
alter table projects add country varchar(100) null;
alter table projects add province varchar(100) null;
alter table projects add city varchar(200) null;
alter table projects add district varchar(200) null;
create index projects_country_province_city_district_index on projects (country, province, city, district);
##2020-05-08
INSERT INTO culture.tags (id, name, code, is_delete, created_at, updated_at) VALUES ('e69fbcc4-18b1-472a-882c-c7fd47bea510', '献礼70周年', '献礼70周年', 0, '2020-05-08 11:26:36.000', '2020-05-08 11:26:38.000')