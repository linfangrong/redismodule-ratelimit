# redismodule-ratelimit

Install and configure redis.

## Usage
From Redis (try running `redis-cli`) use the new `RATELIMIT.ALLOW` command loaded by
the module. It's used like this:
```
RATELIMIT.ALLOW <key> <interval(ns)> <burst>
```
Where `key` is an identifier to rate limit against.

For example:
```
RATELIMIT.ALLOW abc 1000000000 100
                 ▲       ▲      ▲
                 |       |      └─────  100 burst
                 |       └────────────  1 tokens / 1 seconds
                 └────────────────────  key "abc"
```
### Response
The command will respond with integer:
```
127.0.0.1:6379> RATELIMIT.ALLOW abc 1000000000 100
(integer) 1
```
Whether the action was limited:
 * `1` indicates the action is allowed.
 * `0` indicates that the action was limited/blocked.
