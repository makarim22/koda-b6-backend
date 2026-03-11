CREATE TABLE IF NOT EXISTS forgot_password(
    id serial primary key,
    code_otp int not null,
    email varchar(100) unique,
    expired_at timestamp default now() + interval '15 minutes',
    created_at timestamp default now(),
    updated_at timestamp default now()
)