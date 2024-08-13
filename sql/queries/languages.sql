-- name: GetLanguageIDByName :one
SELECT id
FROM languages
WHERE language = $1;

