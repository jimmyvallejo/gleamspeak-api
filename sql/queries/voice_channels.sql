-- name: CreateVoiceChannel :one
INSERT INTO voice_channels (
        id,
        owner_id,
        server_id,
        language_id,
        channel_name,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: JoinVoiceChannel :one
INSERT INTO voice_channel_members (
    user_id,
    channel_id,
    server_id
    )
VALUES ($1, $2, $3)
RETURNING *;
-- name: LeaveVoiceChannelByUser :exec
DELETE FROM voice_channel_members
WHERE user_id = $1;
-- name: GetServerVoiceChannels :many
SELECT 
    vc.id AS channel_id,
    vc.owner_id,
    vc.server_id,
    vc.language_id,
    vc.channel_name,
    vc.last_active,
    vc.is_locked,
    vc.created_at AS channel_created_at,
    vc.updated_at AS channel_updated_at,
    COALESCE(json_agg(
        json_build_object(
            'user_id', vcm.user_id,
            'handle', u.handle
        ) 
    ) FILTER (WHERE vcm.user_id IS NOT NULL), '[]'::json) AS members
FROM 
    voice_channels vc
LEFT JOIN 
    voice_channel_members vcm ON vc.id = vcm.channel_id
LEFT JOIN
    users u ON vcm.user_id = u.id
WHERE 
    vc.server_id = $1
GROUP BY
    vc.id
ORDER BY 
    vc.created_at DESC;