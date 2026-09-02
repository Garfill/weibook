local key = KEYS[1]
local val = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get", key..":cnt"))

if cnt <= 0 then
    --  一直输错
    -- 或者已经用过
    return -1
elseif val == code then
     -- 输入正确
     -- 设置验证码失效
     redis.call("set", cntKey, -1)
     return 0
else
     -- 用户输入错误,次数-1
     redis.call("decr", cntKey)
     return -2
end