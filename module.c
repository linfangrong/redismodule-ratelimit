#include "redismodule.h"
#include "libs/libratelimit.h"

int AllowCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
  if (argc < 4) return RedisModule_WrongArity(ctx);
  RedisModule_AutoMemory(ctx);

  const char *p;
  size_t n;
  p = RedisModule_StringPtrLen(argv[1], &n);
  GoString resource = {p, n};
  long long interval;
  long long burst;
  if (RedisModule_StringToLongLong(argv[2], &interval) != REDISMODULE_OK) {
    return RedisModule_ReplyWithError(ctx, "ERR interval is not an integer or out of range");
  }
  if (RedisModule_StringToLongLong(argv[3], &burst) != REDISMODULE_OK) {
    return RedisModule_ReplyWithError(ctx, "ERR burst is not an integer or out of range");
  }

  GoUint8 allow = Allow(resource, interval, burst);
  if (allow > 0) {
	  RedisModule_ReplyWithLongLong(ctx, 1);
  } else {
	  RedisModule_ReplyWithLongLong(ctx, 0);
  }
  return REDISMODULE_OK;
}

int RedisModule_OnLoad(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
  if (RedisModule_Init(ctx, "ratelimit", 1, REDISMODULE_APIVER_1)
    == REDISMODULE_ERR) {
    return REDISMODULE_ERR;
  }

  if (RedisModule_CreateCommand(ctx, "ratelimit.allow", AllowCommand, "write fast",
    1, 1, 1) == REDISMODULE_ERR) {
    return REDISMODULE_ERR;
  }

  return REDISMODULE_OK;
}
