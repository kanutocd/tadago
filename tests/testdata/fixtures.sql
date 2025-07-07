-- Test data for integration tests
-- Users
INSERT INTO
  users (id, name, email, created_at, updated_at)
VALUES
  (
    '550e8400-e29b-41d4-a716-446655440001',
    'Test User 1',
    'test1@example.com',
    NOW(),
    NOW()
  ),
  (
    '550e8400-e29b-41d4-a716-446655440002',
    'Test User 2',
    'test2@example.com',
    NOW(),
    NOW()
  ),
  (
    '550e8400-e29b-41d4-a716-446655440003',
    'Test User 3',
    'test3@example.com',
    NOW(),
    NOW()
  );

-- Tadas
INSERT INTO
  tadas (
    id,
    name,
    description,
    created_by,
    assigned_to,
    status,
    due_at,
    created_at,
    updated_at
  )
VALUES
  (
    '650e8400-e29b-41d4-a716-446655440001',
    'Test Task 1',
    'Description 1',
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440002',
    'in_progress',
    NOW()+INTERVAL '7 days',
    NOW(),
    NOW()
  ),
  (
    '650e8400-e29b-41d4-a716-446655440002',
    'Test Task 2',
    'Description 2',
    '550e8400-e29b-41d4-a716-446655440002',
    NULL,
    'in_progress',
    NULL,
    NOW(),
    NOW()
  ),
  (
    '650e8400-e29b-41d4-a716-446655440003',
    'Test Task 3',
    'Description 3',
    '550e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440003',
    'completed',
    NOW()-INTERVAL '1 day',
    NOW()-INTERVAL '1 day',
    NOW()
  );

-- Update completed_at for completed task
UPDATE tadas
SET
  completed_at=NOW()-INTERVAL '1 day'
WHERE
  id='650e8400-e29b-41d4-a716-446655440003'
  AND status='completed';
