-- name: GetOngoingGames :many
SELECT games.id, games.fen, games.completed, games.date_started, games.current_state, users1.raw_user_meta_data ->> 'username' as username1, users2.raw_user_meta_data ->> 'username' as username2
FROM games
JOIN player_games j ON games.id = j.game_id
JOIN auth.users users1 ON j.user_id = users1.id AND j.color = 'w'
JOIN player_games j2 ON games.id = j2.game_id AND j2.user_id != j.user_id AND j2.color ='b'
JOIN auth.users users2 ON j2.user_id = users2.id
WHERE (users1.id = $1 OR users2.id = $1) AND games.completed=false;

-- name: GetGame :one
SELECT 
  games.*, 
  user1.raw_user_meta_data ->> 'username' AS player1,
  user2.raw_user_meta_data ->> 'username' AS player2
FROM games 
JOIN player_games AS player_games_1 ON games.id = player_games_1.game_id AND player_games_1.color = 'w' 
JOIN player_games AS player_games_2 ON games.id = player_games_2.game_id AND player_games_2.color = 'b' 
JOIN auth.users AS user1 ON user1.id = player_games_1.user_id
JOIN auth.users AS user2 ON user2.id = player_games_2.user_id
WHERE games.id = $1;

-- name: CreateGame :one
INSERT INTO games (current_state, ruleset, type)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetIdFromUsername :one
SELECT id FROM auth.users
WHERE raw_user_meta_data ->> 'username' = $1;

-- name: GameJunction :exec
INSERT INTO player_games (user_id, game_id, color)
VALUES ($1, $2, $3);

-- name: MakeMove :exec
UPDATE games
SET current_state = $2, history = $3
WHERE id = $1;