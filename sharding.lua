local function shard(txn)
    -- Assigning Shard count
    shard_count = 2

    -- Generic Name for a Read and a Write shard
    write_shard = 'app-write-shard'
    read_only_shard = 'app-read-only-shard'
    
    -- Extracting the Key from the last segment of the URL
    local path = txn.sf:path()
    local segments = {}
    for segment in path:gmatch("[^/]+") do
        table.insert(segments, segment)
    end
    local last_segment = segments[#segments]
    
    -- Calculating Hash of Key
    local shard_number = (string.byte(last_segment, 1)% shard_count)
    local shard_number_string = tostring(shard_number)
    
    local backend
    local method = txn.sf:method()
    if method == "GET" then
        -- determining Read Shard
        backend = read_only_shard .. shard_number_string
        txn:set_var('req.read_shard', backend)
    else
        -- determining Write Shard
        backend = write_shard .. shard_number_string
        txn:set_var('req.write_shard', backend)
    end
    print("Returning ..." .. backend)
    return backend
end


core.register_action('shard', {'http-req'}, shard)