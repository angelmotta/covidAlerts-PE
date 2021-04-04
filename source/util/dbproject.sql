-- DB statments

create database covid;

CREATE TABLE dailycases (
    newcases_date   DATE PRIMARY KEY,
    newcases_amount INT NOT NULL,
    totalcases      INT NOT NULL,
    tsrecord        TIMESTAMPTZ NOT NULL
);

CREATE TABLE dailydeceased (
    deceasedcases_date	DATE PRIMARY KEY,
    newdeceased_amount	INT NOT NULL,
    totaldeceased		INT NOT NULL,
    tsrecord            TIMESTAMPTZ NOT NULL
);

-- Insert new positive cases
INSERT INTO dailycases (newcases_date, newcases_amount, totalcases, tsrecord)
VALUES ('2021-01-31',1000,50,'2021-02-20 12:05:25-07');
-- Insert deceased cases
INSERT INTO dailydeceased (deceasedcases_date, newdeceased_amount, totaldeceased, tsrecord)
VALUES ('2021-01-31',100,5,'2021-02-20 12:17:25');
