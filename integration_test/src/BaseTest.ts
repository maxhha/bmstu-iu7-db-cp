import { getRedisClient } from "./redisClient"
import { getClient } from "./client"
import { UpdateUserDraftFormInput } from "./query"

type FillUserParams = {
  token: string
  email?: string
  oldPassword?: string
  password?: string
  phone?: string
  form?: UpdateUserDraftFormInput
}

export class BaseTest {
  protected client = getClient()
  protected rdb!: Awaited<ReturnType<typeof getRedisClient>>

  async before() {
    this.rdb = await getRedisClient()
  }

  async after() {
    await this.rdb.quit()
  }

  async register() {
    let userId!: string
    let token!: string

    await this.client.Register().then((resp) => {
      expect(resp.status).toBe(200)
      token = resp.data.register.token
      this.client.setToken(resp.data.register.token)
    })

    await this.client.Viewer().then((response) => {
      expect(response.status).toBe(200)
      expect(response.data.viewer).not.toBeNull()
      expect(response.data.viewer).not.toBeUndefined()
      userId = response.data.viewer!.id
    })

    return { userId, token }
  }

  async fillUser({
    token,
    password,
    oldPassword,
    email,
    form,
    phone,
  }: FillUserParams) {
    this.client.setToken(token)

    if (password) {
      console.debug("update password")
      await this.client
        .UpdateUserPassword({
          input: { password, oldPassword },
        })
        .then((response) => expect(response.status).toBe(200))
    }

    if (email) {
      console.debug("set email")
      await this.client
        .RequestSetUserEmail({ input: { email } })
        .then((response) => {
          expect(response.status).toBe(200)
          expect(response.data.requestSetUserEmail).toBe(true)
        })

      const data = await this.rdb.GET(`send-SET_USER_EMAIL-${email}`)
      expect(data).not.toBeNull()
      const parsed = JSON.parse(data!)

      await this.client
        .ApproveSetUserEmail({
          input: { token: parsed.token },
        })
        .then((response) => expect(response.status).toBe(200))
    }

    if (phone) {
      console.debug("set phone")
      await this.client
        .RequestSetUserPhone({ input: { phone } })
        .then((response) => {
          expect(response.status).toBe(200)
          expect(response.data.requestSetUserPhone).toBe(true)
        })

      const data = await this.rdb.GET(`send-SET_USER_PHONE-${phone}`)
      expect(data).not.toBeNull()
      const parsed = JSON.parse(data!)

      await this.client
        .ApproveSetUserPhone({
          input: { token: parsed.token },
        })
        .then((response) => expect(response.status).toBe(200))
    }

    if (form) {
      console.debug("set form")

      await this.client
        .UpdateUserDraftForm({ input: form })
        .then((response) => expect(response.status).toBe(200))
    }
  }
}

export function noop(..._: any) {}

export const randid = () => Math.random().toString(32).slice(2)
