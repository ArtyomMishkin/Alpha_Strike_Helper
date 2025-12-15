-- ==================== SEED DATA FOR CARDS ====================
-- Insert BattleTech Alpha Strike mechs into cards table

INSERT INTO cards (name, model, type, size, role, faction, pv, tech_base, move, tmm, armor, structure, dmg_short, dmg_medium, dmg_long, overheat, abilities, created_at, updated_at) VALUES

-- CLAN MECHS
('Timber Wolf Prime', 'TWF-1', 'Mech', 'Heavy', 'Assault', 'Clans', 45, 'Clan', 6, 1, 19, 16, 6, 5, 4, 0, 'Jump', NOW(), NOW()),
('Nova Cat', 'NOC-A', 'Mech', 'Medium', 'Striker', 'Clans', 32, 'Clan', 8, 2, 12, 10, 4, 4, 3, 1, 'Jump', NOW(), NOW()),
('Mad Cat Mk II', 'MCII-A', 'Mech', 'Heavy', 'Assault', 'Clans', 48, 'Clan', 6, 1, 21, 17, 7, 6, 5, 2, 'Jump', NOW(), NOW()),
('Stormcrow Prime', 'SCR-1', 'Mech', 'Medium', 'Skirmisher', 'Clans', 35, 'Clan', 8, 2, 13, 11, 5, 4, 3, 1, 'Jump', NOW(), NOW()),
('Elemental', 'ELE', 'Mech', 'Light', 'Striker', 'Clans', 18, 'Clan', 10, 3, 5, 4, 2, 2, 1, 0, 'Jump', NOW(), NOW()),

-- INNER SPHERE MECHS
('Locust LCT-1V', 'LCT-1V', 'Mech', 'Light', 'Scout', 'IS', 15, 'IS', 12, 4, 4, 3, 1, 1, 0, 0, 'Jump,Speed Boost', NOW(), NOW()),
('Commando COM-2D', 'COM-2D', 'Mech', 'Light', 'Skirmisher', 'IS', 20, 'IS', 10, 3, 6, 4, 2, 2, 1, 0, 'Jump', NOW(), NOW()),
('Jagermech JM6-S', 'JM6-S', 'Mech', 'Medium', 'Sniper', 'IS', 28, 'IS', 7, 1, 10, 8, 4, 4, 3, 1, 'Target Info Sale', NOW(), NOW()),
('Hunchback HBK-4G', 'HBK-4G', 'Mech', 'Medium', 'Striker', 'IS', 30, 'IS', 8, 2, 11, 9, 5, 4, 3, 0, 'Torso Twist', NOW(), NOW()),
('Enforcer ENF-4R', 'ENF-4R', 'Mech', 'Medium', 'Brawler', 'IS', 32, 'IS', 7, 1, 12, 10, 5, 4, 3, 1, 'Jump', NOW(), NOW()),
('Catapult CPLT-C1', 'CPLT-C1', 'Mech', 'Heavy', 'Assault', 'IS', 42, 'IS', 5, 0, 16, 14, 6, 5, 4, 2, 'Jump', NOW(), NOW()),
('Atlas AS7-D', 'AS7-D', 'Mech', 'Heavy', 'Assault', 'IS', 50, 'IS', 4, -1, 22, 18, 7, 6, 5, 2, 'Torso Twist,Targeting Computer', NOW(), NOW()),
('Awesome AWS-9M', 'AWS-9M', 'Mech', 'Heavy', 'Assault', 'IS', 46, 'IS', 5, 0, 19, 15, 7, 6, 5, 3, 'Jump', NOW(), NOW()),

-- OMNI MECHS
('Loki Prime', 'LKI-1', 'OmniMech', 'Heavy', 'Assault', 'Clans', 44, 'Clan', 6, 1, 17, 15, 6, 5, 4, 1, 'Modular', NOW(), NOW()),
('Summoner Prime', 'SUM-1', 'OmniMech', 'Heavy', 'Brawler', 'Clans', 40, 'Clan', 6, 1, 15, 13, 5, 4, 3, 1, 'Modular', NOW(), NOW()),
('Omni Mech Fighter', 'OMN-X', 'OmniMech', 'Medium', 'Striker', 'Clans', 33, 'Clan', 8, 2, 11, 9, 4, 3, 2, 1, 'Modular', NOW(), NOW()),

-- LIGHT MECHS
('Flea FLE-4', 'FLE-4', 'Mech', 'Light', 'Scout', 'IS', 12, 'IS', 14, 5, 3, 2, 1, 1, 0, 0, 'Jump', NOW(), NOW()),
('Urbanmech UM-R60', 'UM-R60', 'Mech', 'Light', 'Urban Warrior', 'IS', 18, 'IS', 6, 0, 6, 4, 2, 1, 1, 0, 'Stable Firing', NOW(), NOW()),

-- MEDIUM MECHS
('Cicada CDA-2A', 'CDA-2A', 'Mech', 'Light', 'Scout', 'IS', 16, 'IS', 15, 5, 4, 3, 1, 1, 0, 0, 'Jump', NOW(), NOW()),
('Jagermech JM7-A', 'JM7-A', 'Mech', 'Medium', 'Sniper', 'IS', 30, 'IS', 8, 2, 10, 8, 4, 4, 3, 1, 'Targeting Computer', NOW(), NOW()),

-- BATTLE ARMOR
('Elemental Battle Armor', 'BA-EL', 'Battle Armor', 'Light', 'Striker', 'Clans', 5, 'Clan', 10, 2, 2, 2, 1, 1, 0, 0, 'Squad', NOW(), NOW()),
('IS Standard BA', 'BA-STD', 'Battle Armor', 'Light', 'Trooper', 'IS', 4, 'IS', 8, 2, 2, 1, 1, 1, 0, 0, 'Squad', NOW(), NOW());

-- ==================== VERIFY INSERT ====================
-- SELECT COUNT(*) as total_cards FROM cards;
-- SELECT name, role, pv FROM cards ORDER BY pv DESC;
