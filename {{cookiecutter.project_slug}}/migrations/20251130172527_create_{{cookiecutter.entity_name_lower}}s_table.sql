-- +goose Up
-- +goose StatementBegin
CREATE TABLE {{cookiecutter.entity_name_lower}}s (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_{{cookiecutter.entity_name_lower}}s_name ON {{cookiecutter.entity_name_lower}}s(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS {{cookiecutter.entity_name_lower}}s;
-- +goose StatementEnd
