# redismodule-ratelimit

Install and configure redis.

## Usage
From Redis (try running `redis-cli`) use the new `RATELIMIT.ALLOW` command loaded by
the module. It's used like this:
```
RATELIMIT.ALLOW <key> <interval> <burst>
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

