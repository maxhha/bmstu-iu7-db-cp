import { createClient } from "redis"

export async function getRedisClient() {
  if (process.env.REDIS_CONNECTION_URL === undefined) {
    throw new Error("REDIS_CONNECTION_URL environment variable is undefined")
  }

  const client = createClient({ url: process.env.REDIS_CONNECTION_URL })

  client.on("error", (err) => console.error("redis error:", err))

  await client.connect()
  return client
}
