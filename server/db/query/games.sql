-- -- name: GetGames :many
-- SELECT games.*, users.raw_user_meta_data -> 'username' as username
-- FROM games
-- INNER JOIN player_games ON games.id = player_games.game_id
-- INNER JOIN auth.users ON player_games.user_id = users.id
-- WHERE users.id = $1;

-- name: GetGames :many
SELECT games.*, users1.raw_user_meta_data -> 'username' as username1, users2.raw_user_meta_data -> 'username' as username2
FROM games
JOIN player_games j ON games.id = j.game_id
JOIN auth.users users1 ON j.user_id = users1.id
JOIN player_games j2 ON games.id = j2.game_id AND j2.user_id != j.user_id
JOIN auth.users users2 ON j2.user_id = users2.id
WHERE users1.id = $1 AND j.color ='w' OR users2.id = $1 AND j.color ='b';

-- -- name: GetGames :one
-- SELECT * FROM games;