-- name: GetOngoingGames :many
SELECT 
    games.id, 
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

-- name: GetGame :one
SELECT 
    games.id,
    games.fen,
    games.history,
    games.completed,
    games.date_started,
    games.date_finished,
    games.current_state,
    games.ruleset,
    games.type,
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
    games.id = $1;

-- name: GetGameWithUndo :one
SELECT 
    games.id,
    games.fen,
    games.history,
    games.completed,
    games.date_started,
    games.date_finished,
    games.current_state,
    games.ruleset,
    games.type,
    games.user_1,
    games.user_2,
    user1.username AS player1,
    user2.username AS player2,
    CASE
        WHEN COUNT(undo_request.game_id) > 0 THEN
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
    undo_request ON undo_request.game_id = games.id
LEFT JOIN 
    profiles AS sender ON sender.id = undo_request.sender_id
LEFT JOIN 
    profiles AS receiver ON receiver.id = undo_request.receiver_id
WHERE
    games.id = $1
GROUP BY
    games.id, user1.username, user2.username;

-- name: CreateGame :one
INSERT INTO games (current_state, ruleset, type, user_1, user_2)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetIdFromUsername :one
SELECT id FROM profiles
WHERE profiles.username = $1;

-- name: GetUsernameFromId :one
SELECT username FROM profiles
WHERE profiles.id = $1;

-- name: GetOnboarding :one
SELECT is_username_onboard_complete FROM profiles
WHERE profiles.id = $1;

-- name: UpdateOnboarding :exec
UPDATE profiles
SET is_username_onboard_complete = true
WHERE profiles.id = $1;

-- name: ChangeUsername :exec
UPDATE profiles
SET username = $1
WHERE profiles.id = $2;

-- name: MakeMove :exec
UPDATE games
SET current_state = $2, history = $3
WHERE id = $1;

-- name: CreateUndo :one
INSERT INTO undo_request (game_id, sender_id, receiver_id)
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
            games.id = $1
    )
)
RETURNING receiver_id;

-- name: RemoveUndo :exec
DELETE FROM undo_request
WHERE sender_id = $1 AND game_id = $2;

-- name: ChangeUndo :one
UPDATE undo_request
SET status = $1
WHERE receiver_id = $2 AND game_id = $3
RETURNING undo_request.sender_id;

-- name: CreateRoom :exec
INSERT INTO public.room_list (host_id, description, rules, type, color)
VALUES ($1, $2, $3, $4, $5);

-- name: GetRoomList :many
SELECT
    room_list.*,
    profiles.username AS host
FROM
    public.room_list
JOIN
    public.profiles ON room_list.host_id = profiles.id;

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