-- DB statments

create database covid;

CREATE TABLE dailycases (
    newcases_date   DATE PRIMARY KEY,
    newcases_amount INT NOT NULL
);

CREATE TABLE dailydeceased (
    deceasedcases_date	DATE PRIMARY KEY,
    newdeceased_amount	INT NOT NULL,
    totaldeceased		INT NOT NULL
);

-- Insert
INSERT INTO dailycases (newcases_date, newcases_amount) VALUES ('2021-01-31', 2710)