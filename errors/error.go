package errors

const Nil = RedisError("redis: nil")

type RedisError string

func (e RedisError) Error() string { return string(e) }
