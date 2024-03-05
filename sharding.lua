local function shard(txn)
    shard_count = 2
    write_shard = 'app-write-shard'
    read_only_shard = 'app-read-only-shard'
    
    local path = txn.sf:path()
    local segments = {}
    for segment in path:gmatch("[^/]+") do
        table.insert(segments, segment)
    end

    local method = txn.sf:method()
    local last_segment = segments[#segments]
    
    local shard_number = (string.byte(last_segment, 1)% shard_count)
    local shard_number_string = tostring(shard_number)
    local backend
    print("method", method)

    if method == "GET" then
        backend = read_only_shard .. shard_number_string
        txn:set_var('req.read_shard', backend)
    else
        backend = write_shard .. shard_number_string
        txn:set_var('req.write_shard', backend)
    end
    print("Returning ..." .. backend)
    return backend
end


core.register_action('shard', {'http-req'}, shard)