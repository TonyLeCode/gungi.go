-- name: GetOngoingGames :many
SELECT 
    games.public_id, 
    games.fen, 
    games.completed, 
    games.date_started, 
    games.current_state, 
    users1.username as username1, 
    users2.username as username2
FROM 
    games
JOIN 
    profiles as users1 ON games.user_1 = users1.id
JOIN 
    profiles as users2 ON games.user_2 = users2.id
WHERE
    (users1.id = $1 OR users2.id = $1) AND games.completed=false;

-- name: GetOverview :many
SELECT 
    games.public_id, 
    games.fen, 
    games.completed, 
    games.date_started, 
    games.date_finished,
    games.current_state,
    games.result,
    games.type,
    games.ruleset,
    users1.username as username1, 
    users2.username as username2
FROM
    games 
JOIN 
    profiles as users1 ON games.user_1 = users1.id
JOIN 
    profiles as users2 ON games.user_2 = users2.id
WHERE
    users1.id = $1 OR users2.id = $1 AND games.completed=false;

-- name: GetCompletedGames :many
SELECT 
    games.public_id, 
    games.fen, 
    games.completed, 
    games.date_started, 
    games.date_finished,
    games.current_state,
    games.result,
    games.type,
    games.ruleset,
    users1.username as username1, 
    users2.username as username2
FROM
    games 
JOIN 
    profiles as users1 ON games.user_1 = users1.id
JOIN 
    profiles as users2 ON games.user_2 = users2.id
WHERE
    users1.id = $1 OR users2.id = $1 AND games.completed=true
ORDER BY games.date_finished DESC
OFFSET $2 ROWS FETCH NEXT 10 ROWS ONLY;

-- name: GetCompletedGamesCount :one
SELECT COUNT(*) FROM games
WHERE games.user_1 = $1 OR games.user_2 = $1
AND games.completed=true;

-- name: GetGame :one
SELECT 
    games.public_id,
    games.fen,
    games.history,
    games.completed,
    games.date_started,
    games.date_finished,
    games.current_state,
    games.ruleset,
    games.type,
    games.result,
    games.user_1,
    games.user_2,
    user1.username AS player1,
    user2.username AS player2
FROM 
    games
JOIN 
    profiles AS user1 ON user1.id = games.user_1
JOIN 
    profiles AS user2 ON user2.id = games.user_2
WHERE
    games.public_id = $1;

-- name: GetGameWithUndo :one
SELECT 
    games.public_id,
    games.fen,
    games.history,
    games.completed,
    games.date_started,
    games.date_finished,
    games.current_state,
    games.ruleset,
    games.result,
    games.type,
    games.user_1,
    games.user_2,
    user1.username AS player1,
    user2.username AS player2,
    CASE
        WHEN COUNT(undo_request.game_public_id) > 0 THEN
            json_agg(
                json_build_object(
                    'sender_username', sender.username,
                    'receiver_username', receiver.username,
                    'status', undo_request.status
                )
            )
        ELSE
            '[]'::json
    END AS undo_requests
FROM 
    games
JOIN 
    profiles AS user1 ON user1.id = games.user_1
JOIN 
    profiles AS user2 ON user2.id = games.user_2
LEFT JOIN
    undo_request ON undo_request.game_public_id = games.public_id
LEFT JOIN 
    profiles AS sender ON sender.id = undo_request.sender_id
LEFT JOIN 
    profiles AS receiver ON receiver.id = undo_request.receiver_id
WHERE
    games.public_id = $1
GROUP BY
    games.id, user1.username, user2.username;

-- name: CreateGame :exec
INSERT INTO games (public_id, current_state, ruleset, type, user_1, user_2)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetIdFromUsername :one
SELECT id FROM profiles
WHERE profiles.username = $1;

-- name: GetUsernameFromId :one
SELECT username FROM profiles
WHERE profiles.id = $1;

-- name: MakeMove :exec
UPDATE games
SET current_state = $2, history = $3
WHERE games.public_id = $1;

-- name: ChangeGameResult :exec
UPDATE games
SET completed = $2, result = $3, date_finished = now()
WHERE games.public_id = $1;

-- name: ResignGame :one
WITH updated_game AS (
    UPDATE games
    SET 
        completed = true,
        result = CASE 
            WHEN user_1 = $2 THEN 'b/r'
            WHEN user_2 = $2 THEN 'w/r'
            ELSE result
        END,
        date_finished = now()
    WHERE games.public_id = $1 AND (user_1 = $2 OR user_2 = $2)
    RETURNING games.public_id, fen, history, completed, date_started, date_finished, current_state, ruleset, type, result, user_1, user_2
)
SELECT 
    ug.public_id,
    ug.fen,
    ug.history,
    ug.completed,
    ug.date_started,
    ug.date_finished,
    ug.current_state,
    ug.ruleset,
    ug.type,
    ug.result,
    ug.user_1,
    ug.user_2,
    user1.username AS player1,
    user2.username AS player2
FROM 
    updated_game ug
JOIN 
    profiles AS user1 ON user1.id = ug.user_1
JOIN 
    profiles AS user2 ON user2.id = ug.user_2;


-- name: CreateUndo :one
INSERT INTO undo_request (game_public_id, sender_id, receiver_id)
VALUES (
    $1,
    $2,
    (
        SELECT
            CASE
                WHEN games.user_1 = $2 THEN games.user_2
                ELSE games.user_1
            END
        FROM
            games
        WHERE
            games.public_id = $1 AND games.completed = false
    )
)
RETURNING receiver_id;

-- name: RemoveUndo :exec
DELETE FROM undo_request
WHERE sender_id = $1 AND game_public_id = $2;

-- name: ChangeUndo :one
UPDATE undo_request
SET status = $1
WHERE receiver_id = $2 AND game_public_id = $3
RETURNING undo_request.sender_id;

-- name: CreateRoom :exec
INSERT INTO public.room_list (host_id, description, rules, type, color)
VALUES ($1, $2, $3, $4, $5);

-- name: GetRoomList :many
SELECT
    room_list.id,
    room_list.description,
    room_list.rules,
    room_list.type,
    room_list.color,
    profiles.username AS host
FROM
    public.room_list
JOIN
    public.profiles ON room_list.host_id = profiles.id;

-- name: DeleteRoomSafe :one
WITH deleted_room AS (
    DELETE FROM public.room_list
    WHERE room_list.id = $1 AND room_list.host_id = $2
    AND EXISTS (SELECT FROM profiles WHERE profiles.id = room_list.host_id)
    RETURNING *
)
SELECT deleted_room.*, profiles.username as host
FROM deleted_room
JOIN profiles ON deleted_room.host_id = profiles.id;

-- name: DeleteRoom :one
WITH deleted_room AS (
    DELETE FROM public.room_list
    WHERE room_list.id = $1
    AND EXISTS (SELECT FROM profiles WHERE profiles.id = room_list.host_id)
    RETURNING *
)
SELECT deleted_room.*, profiles.username as host
FROM deleted_room
JOIN profiles ON deleted_room.host_id = profiles.id;