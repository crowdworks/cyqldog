-- SQLs for setup of development environment

-- Create test table
DROP TABLE IF EXISTS table1;
CREATE TABLE table1 (
    tag1 char(8),
    val1 integer,
    tag2 varchar(8),
    val2 real
);

-- Insert test data
INSERT INTO table1 (tag1, val1, tag2, val2) VALUES ('hoge1', 1, 'fuga1', 0.1);
INSERT INTO table1 (tag1, val1, tag2, val2) VALUES ('hoge1', 2, 'fuga2', 0.2);
INSERT INTO table1 (tag1, val1, tag2, val2) VALUES ('hoge3', 3, 'fuga3', 0.3);
