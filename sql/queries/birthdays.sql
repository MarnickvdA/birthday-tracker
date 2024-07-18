-- name: ScheduleBirthdayNotifications :one
WITH inserted_rows AS (
    INSERT INTO birthday_notifications (person_id, scheduled_at, state)
    SELECT id, TO_CHAR(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 'scheduled'
    FROM persons
    WHERE TO_CHAR(TO_DATE(birth_date, 'YYYY-MM-DD'), 'MM-DD') = TO_CHAR(CURRENT_DATE + INTERVAL '1 day', 'MM-DD')
    ON CONFLICT (person_id) DO NOTHING
    RETURNING *
)
SELECT COUNT(*) FROM inserted_rows;

-- name: GetScheduledBirthdayNotificationsForToday :many
SELECT 
    bn.person_id,
    p.name, 
    p.birth_date,
    EXTRACT(YEAR FROM AGE(TO_DATE(p.birth_date, 'YYYY-MM-DD'))) AS age
FROM 
    birthday_notifications bn
JOIN 
    persons p ON bn.person_id = p.id
WHERE
    TO_CHAR(TO_DATE(bn.scheduled_at, 'YYYY-MM-DD'), 'MM-DD') = TO_CHAR(CURRENT_DATE, 'MM-DD')
    AND bn.state = 'scheduled';

-- name: UpdateBirthdayNotificationsStateForToday :exec
UPDATE birthday_notifications
SET state = $1
WHERE TO_DATE(scheduled_at, 'YYYY-MM-DD') = CURRENT_DATE;

-- name: DeletePastBirthdayNotifications :one
WITH deleted_rows AS (
    DELETE FROM birthday_notifications
    WHERE TO_DATE(scheduled_at, 'YYYY-MM-DD') < CURRENT_DATE
    RETURNING *
)
SELECT COUNT(*) FROM deleted_rows;