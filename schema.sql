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
SELECT inr.a, avg(inr.price) as price, inr.b, sum(inr.quantity) as market
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
