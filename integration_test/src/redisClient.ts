import { createClient } from "redis"

let client: Awaited<ReturnType<typeof createClient>>

export async function getRedisClient() {
  if (client) return client

  if (process.env.REDIS_CONNECTION_URL === undefined) {
    throw new Error("REDIS_CONNECTION_URL environment variable is undefined")
  }

  client = createClient({ url: process.env.REDIS_CONNECTION_URL })

  client.on("error", (err) => console.error("redis error:", err))

  await client.connect()
  return client
}
