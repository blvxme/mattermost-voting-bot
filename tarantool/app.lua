local user = os.getenv("TARANTOOL_USER")
local password = os.getenv("TARANTOOL_PASSWORD")

if not user then
    error("Environment variable TARANTOOL_USER not set")
end

if not password then
    error("Environment variable TARANTOOL_PASSWORD not set")
end

box.cfg{ listen="0.0.0.0:3301" }
box.schema.user.create(user, { password = password, if_not_exists = true })
box.schema.user.grant(user, 'super', nil, nil, { if_not_exists = true })
require('msgpack').cfg{ encode_invalid_as_nil = true }

box.schema.space.create("votings", { if_not_exists = true })
box.space.votings:format({
    { name = "id",         type = "string"                    }, -- ID голосования
    { name = "creator_id", type = "string"                    }, -- ID пользователя, который создал голосование
    { name = "is_ended",   type = "boolean"                   }, -- Закончено ли голосование
    { name = "title",      type = "string"                    }, -- Заголовок голосования
    { name = "options",    type = "map"                       }, -- Варианты ответа (вариант: количество голосов)
    { name = "users",      type = "array", is_nullable = true }  -- Проголосовавшие пользователи
})

box.space.votings:create_index("primary", { parts = { "id" }, if_not_exists = true })
