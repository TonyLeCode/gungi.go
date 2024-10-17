-- name: CreateUsername :exec
INSERT INTO public.profiles (id, username)
VALUES ($1, $2);