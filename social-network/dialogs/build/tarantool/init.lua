os.execute('mkdir -p /var/lib/tarantool/memtx /var/lib/tarantool/wal')

box.cfg {
    listen = 3301,

    memtx_dir = '/var/lib/tarantool/memtx',
    wal_dir   = '/var/lib/tarantool/wal',
    log       = '/var/lib/tarantool/tarantool.log',

    memtx_memory = 1024 * 1024 * 1024, -- 1GB

    checkpoint_interval = 60,
    checkpoint_count = 5,

    wal_mode = 'write',
}

if not box.space.dialog_messages then
    local s = box.schema.space.create('dialog_messages', {
        if_not_exists = true,
        format = {
            { name = 'id', type = 'unsigned' },
            { name = 'dialog_id', type = 'string' },
            { name = 'from', type = 'uuid' },
            { name = 'to', type = 'uuid' },
            { name = 'text', type = 'string' },
            { name = 'created_at', type = 'datetime' },
        }
    })

    s:create_index('primary', { parts = { { field = 'id', type = 'unsigned' } } })
    s:create_index('dialog_id_idx', { parts = { { field = 'dialog_id', type = 'string' } }, unique = false })
end

if not box.sequence.message_id then
    box.schema.sequence.create('dialog_message_id', { if_not_exists = true })
end

function add_message(dialog_id, from, to, text, created_at)
    local new_id = box.sequence.dialog_message_id:next()

    box.space.dialog_messages:insert{ new_id, dialog_id, from, to, text, created_at }
    return { status = "ok", id = new_id }
end

function get_messages_by_dialog(dialog_id)
    local tuples = box.space.dialog_messages.index.dialog_id_idx:select(dialog_id)

    local messages = {}
    for _, t in ipairs(tuples) do
        table.insert(messages, {
            id         = t.id,
            dialog_id  = t.dialog_id,
            from       = t.from,
            to         = t.to,
            text       = t.text,
            created_at = t.created_at
        })
    end

    table.sort(messages, function(a, b)
        return a.created_at < b.created_at
    end)

    return messages
end




