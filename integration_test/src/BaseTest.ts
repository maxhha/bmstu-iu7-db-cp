import { getRedisClient } from "./redisClient"
import { getClient } from "./client"

export class BaseTest {
  protected client = getClient()
  protected rdb!: Awaited<ReturnType<typeof getRedisClient>>
  before = async () => {
    this.rdb = await getRedisClient()
  }

  after = async () => {
    await this.rdb.quit()
  }
}

export function noop(..._: any) {}
