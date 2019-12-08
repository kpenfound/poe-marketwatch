CREATE TABLE uniques (
  id varchar(127) primary key,
  name varchar(127),
  corrupted bool,
  original_price decimal,
  original_price_currency varchar(31),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX uniques_idx_created_at ON uniques (created_at DESC NULLS LAST);

CREATE TABLE currency (
  id varchar(127) primary key,
  type varchar(127),
  original_price decimal,
  original_price_currency varchar(31),
  original_quantity integer,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX currency_idx_created_at ON currency (created_at DESC NULLS LAST);

CREATE TABLE divination_cards (
  id varchar(127) primary key,
  name varchar(127),
  mods varchar(263),
  max_stack_size integer,
  original_price decimal,
  original_price_currency varchar(31),
  original_quantity integer,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX divination_cards_idx_created_at ON divination_cards (created_at DESC NULLS LAST);

CREATE TABLE api_pages (
  id varchar(127) primary key,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX api_pages_idx_created_at ON api_pages (created_at DESC NULLS LAST);

CREATE TABLE currency_dictionary (
  id varchar(15),
  name varchar(63)
);
INSERT INTO currency_dictionary (id, name) VALUES 
('chaos', 'Chaos Orb'),
('exa', 'Exalted Orb'),
('mir', 'Mirror of Kalandra'),
('fuse', 'Orb of Fusing'),
('chisel', 'Cartographer''s Chisel'),
('alt', 'Alteration Orb'),
('gcp', 'Gemcutter''s Prism'),
('divine', 'Divine Orb'),
('regal', 'Regal Orb');

CREATE VIEW currency_conversion AS
SELECT 
  inr.a as listed, 
  avg(inr.price) as average_price,
  mode() within group (order by price asc) as mode_price,
  inr.b as listed_for, 
  sum(inr.quantity) as market
FROM (
  SELECT cd.id as a, (original_price / original_quantity) as price, original_price_currency as b, original_quantity as quantity
    FROM currency c
    INNER JOIN currency_dictionary cd ON cd."name" = c."type"
    WHERE created_at > now() - interval '3 days'
    AND original_quantity > 0
    AND original_price_currency IN ('chaos', 'exa', 'mir', 'fuse', 'chisel', 'alt', 'gcp', 'divine', 'regal')
) inr
WHERE inr.a != inr.b
GROUP BY inr.a, inr.b
ORDER BY inr.a, inr.b
;

CREATE VIEW equal_currency_rates AS
SELECT
  distinct listed as listed,
  1 as average_price,
  1 as mode_price,
  listed as listed_for,
  10000 as market
FROM currency_conversion
;

CREATE VIEW currency_rates AS
SELECT * FROM currency_conversion WHERE mode_price > 0
UNION
SELECT * FROM equal_currency_rates
;

CREATE MATERIALIZED VIEW mat_chaos_currency_rates AS
SELECT listed as currency,
  mode_price as rate
  FROM currency_rates
  WHERE listed_for='chaos'
;

CREATE VIEW unique_prices AS
SELECT
  "name",
  corrupted,
  avg(chaos_price) as average_price,
  mode() WITHIN GROUP (order by chaos_price) as mode_price,
  percentile_disc(0.25) WITHIN GROUP (order by chaos_price) as percentile_25,
  count(*) as "availability"
  FROM (
    SELECT "name",
    corrupted,
    (original_price * COALESCE((
      SELECT mccr.rate 
      FROM mat_chaos_currency_rates mccr
      WHERE mccr.currency = original_price_currency
      ), 0)
    ) as chaos_price
    FROM uniques
    WHERE created_at > now() - interval '3 days'
  ) chaos_priced_uniques
  WHERE chaos_price >= 1
  GROUP BY "name", "corrupted"
  ORDER BY mode_price desc
;

CREATE VIEW divination_prices AS
SELECT
  "name",
  avg(chaos_price) as average_price,
  mode() WITHIN GROUP (order by chaos_price) as mode_price,
  percentile_disc(0.25) WITHIN GROUP (order by chaos_price) as percentile_25,
  sum(original_quantity) as "availability"
  FROM (
    SELECT "name",
    (original_price * COALESCE((
      SELECT mccr.rate 
      FROM mat_chaos_currency_rates mccr
      WHERE mccr.currency = original_price_currency
      ), 0)    
    ) / original_quantity as chaos_price,
    original_quantity
    FROM divination_cards
    WHERE created_at > now() - interval '3 days'
  ) chaos_priced_uniques
  WHERE chaos_price >= 1
  GROUP BY "name"
  ORDER BY mode_price desc
;

CREATE MATERIALIZED VIEW mat_div_unique_pairs AS
SELECT u."name" as unique_name, d."name" as card_name, d.max_stack_size, d.mods
FROM (SELECT distinct "name" from uniques) u
INNER JOIN (SELECT distinct "name", max_stack_size, mods FROM divination_cards) d ON d."mods" LIKE concat('%<uniqueitem>{', u."name", '}%')
;

CREATE VIEW divination_card_profits AS
SELECT mdup.unique_name, 
  mdup.card_name, 
  round((up.percentile_25 / mdup.max_stack_size) - dp.percentile_25, 1) as profit,
  round(dp.percentile_25, 1) as card_price, 
  round(up.percentile_25, 1) as unique_price,
  mdup.max_stack_size as stack_size
FROM mat_div_unique_pairs mdup
INNER JOIN divination_prices dp ON mdup.card_name = dp."name"
INNER JOIN unique_prices up ON mdup.unique_name = up."name"
ORDER BY profit desc
;
