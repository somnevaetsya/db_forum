create extension if not exists citext;

create unlogged table if not exists users
(
    nickname citext collate "C" not null unique primary key,
    fullname citext             not null,
    about    text,
    email    citext             not null unique
);

create unlogged table if not exists forums
(
    title   text   not null,
    user_   citext not null references users (nickname),
    slug    citext not null primary key,
    posts   int default 0,
    threads int default 0
);

create unlogged table if not exists threads
(
    id      serial not null primary key,
    title   text      not null,
    author  citext    not null references users (nickname),
    forum   citext    not null references forums (slug),
    message text      not null,
    votes   int                      default 0,
    slug    citext,
    created timestamp with time zone default now()
);

create unlogged table if not exists posts
(
    id        bigserial not null unique primary key,
    parent    int references posts (id),
    author    citext    not null references users (nickname),
    message   text      not null,
    is_edited bool                     default false,
    forum     citext    not null references forums (slug),
    thread    int       not null references threads (id),
    created   timestamp with time zone default now(),
    path      bigint[]                 default array []::integer[]
);

create unlogged table if not exists votes
(
    nickname citext not null references users (nickname),
    thread   int    not null references threads (id),
    voice    int    not null,
    constraint user_thread_key unique (nickname, thread)
);

create unlogged table if not exists user_forum
(
    nickname citext collate "C" not null references users (nickname),
    forum    citext             not null references forums (slug),
    constraint user_forum_key unique (nickname, forum)
);

-- Триггеры и процедуры
create or replace function insert_votes_proc()
    returns trigger as
$$
begin
    update threads set votes = threads.votes + new.voice where id = new.thread;
    return new;
end;
$$ language plpgsql;

create trigger insert_votes
    after insert
    on votes
    for each row
execute procedure insert_votes_proc();


create or replace function update_votes_proc()
    returns trigger as
$$
begin
    update threads set votes = threads.votes + NEW.voice - OLD.voice where id = NEW.thread;
    return NEW;
end;
$$ language plpgsql;

create trigger update_votes
    after update
    on votes
    for each row
execute procedure update_votes_proc();


create or replace function insert_post_before_proc()
    returns trigger as
$$
declare
    parent_post_id posts.id%type := 0;
begin
    new.path = (select path from posts where id = new.parent) || new.id;
    return new;
end;
$$ language plpgsql;

create trigger insert_post_before
    before insert
    on posts
    for each row
execute procedure insert_post_before_proc();

create or replace function insert_post_after_proc()
    returns trigger as
$$
begin
    update forums set posts = forums.posts + 1 where slug = new.forum;
    return new;
end;
$$ language plpgsql;

create trigger insert_post_after
    after insert
    on posts
    for each row
execute procedure insert_post_after_proc();


create or replace function insert_threads_proc()
    returns trigger as
$$
begin
    update forums set threads = forums.threads + 1 where slug = NEW.forum;
    RETURN NEW;
end;
$$ language plpgsql;

create trigger insert_threads
    after insert
    on threads
    for each row
execute procedure insert_threads_proc();


create or replace function add_user()
    returns trigger as
$$
begin
    insert into user_forum (nickname, forum) values (NEW.author, NEW.forum) on conflict do nothing;
    return NEW;
end;
$$ language plpgsql;

create trigger insert_new_thread
    after insert
    on threads
    for each row
execute procedure add_user();

create trigger insert_new_post
    after insert
    ON posts
    for each row
execute procedure add_user();

CREATE INDEX IF NOT EXISTS users_idx on users (nickname, email) include (about, fullname);
create index if not exists users_nickname_hash on users using hash (nickname);
create index if not exists user_forum_all on user_forum (forum, nickname);

create index if not exists forums_slug on forums using hash (slug); --getforum

create index if not exists threads_created on threads using hash (created); --sortthreads
create index if not exists threads_slug on threads using hash (slug); --getthread
create index if not exists thread__forum__hash ON threads using hash (forum);
create index if not exists threads_forum_created on threads (forum, created); --getthreads
create index if not exists threads_id ON threads USING hash (id); --getthread

create index if not exists posts_id_hash on posts using hash (id);
create index if not exists posts_threads on posts using hash (thread);
create index if not exists posts_id_thread_parent_path1 on posts ((path[1]), path); --parenttree
create index if not exists posts_thread_past on posts (thread, path); --flat,tree
create index if not exists posts_thread_id ON posts (thread, id); -- Sort flat
create index if not exists posts_thread_path_idx ON posts (thread, path); -- Sort tree
create index if not exists posts__thread_path_1_idx ON posts (thread, (path[1])); -- Sort parent tree

create unique index if not exists votes_nickname on votes (thread, nickname); --vote