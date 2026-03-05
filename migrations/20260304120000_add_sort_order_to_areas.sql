-- +goose Up
-- Add sort_order column to areas table

ALTER TABLE areas ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0;

-- Update existing areas with sequential sort_order
UPDATE areas SET sort_order = (
    SELECT row_num - 1 FROM (
        SELECT id, ROW_NUMBER() OVER (ORDER BY created_at) - 1 as row_num
        FROM areas
    ) ranked
    WHERE ranked.id = areas.id
);

-- +goose Down
-- Remove sort_order column from areas table

ALTER TABLE areas DROP COLUMN sort_order;
