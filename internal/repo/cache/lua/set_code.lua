-- 验证码在redis的key
-- phone_code:login:15900000000
local key = KEYS[1]
-- 验证次数
-- phone_code:login:15900000000:cnt
local cntKey = key..":cnt"
-- 验证码
local val = ARGV[1]
-- 该键剩余生存时间
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- key 存在，但没有过期时间，有人删了过期时间
    -- 系统错误，设置了key，但没有过期时间
    return -2

    -- key 不存在
    -- 600 - 60 = 540 九分钟
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送频繁
    return -1
end



