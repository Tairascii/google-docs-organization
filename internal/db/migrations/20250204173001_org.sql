-- +goose Up
-- +goose StatementBegin
create table if not exists org (
    id uuid primary key not null default gen_random_uuid(),
    title varchar(64) not null,
    org_picture_url text,
    created_at timestamptz default now()
);

create table if not exists org_members (
    org_id uuid not null,
    user_id uuid not null,
    role text check (role in ('admin', 'member', 'reporter')) not null,
    primary key (org_id, user_id),
    constraint fk_org_id foreign key (org_id) references "org"(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
