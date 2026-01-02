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
            { name = 'state', type = 'string' },
            { name = 'created_at', type = 'unsigned' },
        }
    })

    s:create_index('primary', { parts = { { field = 'id', type = 'unsigned' } } })
    s:create_index('dialog_id_idx', { parts = { { field = 'dialog_id', type = 'string' } }, unique = false })
    s:create_index('dialog_id_id_idx', {
        parts = {
            { field = 'dialog_id', type = 'string' },
            { field = 'id', type = 'unsigned' },
        },
        unique = false,
    })
end

if not box.sequence.message_id then
    box.schema.sequence.create('dialog_message_id', { if_not_exists = true })
end

function add_message(dialog_id, from, to, text, state, created_at)
    local new_id = box.sequence.dialog_message_id:next()

    box.space.dialog_messages:insert{ new_id, dialog_id, from, to, text, state, created_at }
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
            state      = t.state,
            created_at = t.created_at
        })
    end

    table.sort(messages, function(a, b)
        return a.created_at < b.created_at
    end)

    return messages
end

function mark_messages_as_reading(dialog_id, reader_uuid, last_message_id, batch_size)
    batch_size = batch_size or 100

    local space = box.space.dialog_messages
    local index = space.index.dialog_id_id_idx

    local updated_ids = {}
    local batch = {}

    for _, t in index:pairs({dialog_id, 0}, { iterator = box.index.GE }) do
        if t.dialog_id ~= dialog_id then
            break
        end

        if t.id > last_message_id then
            break
        end

        if t.state == 'sent' and t.from ~= reader_uuid then
            table.insert(batch, t.id)
        end

        if #batch >= batch_size then
            box.begin()
            for _, id in ipairs(batch) do
                local cur = space:get(id)
                if cur and cur.state == 'sent' and cur.from ~= reader_uuid then
                    space:update(id, {{'=', 'state', 'reading'}})
                    table.insert(updated_ids, id)
                end
            end
            box.commit()
            batch = {}
        end
    end

    if #batch > 0 then
        box.begin()
        for _, id in ipairs(batch) do
            local cur = space:get(id)
            if cur and cur.state == 'sent' and cur.from ~= reader_uuid then
                space:update(id, {{'=', 'state', 'reading'}})
                table.insert(updated_ids, id)
            end
        end
        box.commit()
    end

    return updated_ids
end

function update_message_state_from(dialog_id, message_id, from_state, to_state)
    local space = box.space.dialog_messages

    box.begin()

    local msg = space:get(message_id)
    if not msg or msg.dialog_id ~= dialog_id then
        box.commit()
        return { updated = false, reason = "not_found" }
    end

    if msg.state ~= from_state then
        box.commit()
        return { updated = false, reason = "wrong_state" }
    end

    space:update(message_id, {{'=', 'state', to_state}})

    box.commit()
    return { updated = true }
end

function update_messages_state_from(
    dialog_id,
    message_ids,
    from_state,
    to_state,
    batch_size
)
    batch_size = batch_size or 100

    local space = box.space.dialog_messages
    local updated_ids = {}

    local batch = {}

    for _, message_id in ipairs(message_ids) do
        table.insert(batch, message_id)

        if #batch >= batch_size then
            box.begin()
            for _, id in ipairs(batch) do
                local msg = space:get(id)
                if msg
                   and msg.dialog_id == dialog_id
                   and msg.state == from_state then

                    space:update(id, {{'=', 'state', to_state}})
                    table.insert(updated_ids, id)
                end
            end
            box.commit()
            batch = {}
        end
    end

    if #batch > 0 then
        box.begin()
        for _, id in ipairs(batch) do
            local msg = space:get(id)
            if msg
               and msg.dialog_id == dialog_id
               and msg.state == from_state then

                space:update(id, {{'=', 'state', to_state}})
                table.insert(updated_ids, id)
            end
        end
        box.commit()
    end

    return updated_ids
end
