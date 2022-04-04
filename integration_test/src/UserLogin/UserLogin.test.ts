import { getClient } from "../client"
import { getRedisClient } from "../redisClient"

it("should create user and login", async () => {
  const password = "HelloWorld-password"
  const email = "test-email@email.com"
  const client = getClient()
  const rdb = await getRedisClient()

  {
    const response = await client.Register()
    expect(response.status).toBe(200)
    client.setToken(response.data.register.token)
  }

  let userId: string
  {
    const response = await client.Viewer()
    expect(response.status).toBe(200)
    expect(response.data.viewer).not.toBeNull()
    expect(response.data.viewer).not.toBeUndefined()
    userId = response.data.viewer!.id
  }

  {
    const response = await client.UpdateUserPassword({ input: { password } })
    expect(response.status).toBe(200)
  }

  {
    const response = await client.RequestSetUserEmail({ input: { email } })
    expect(response.status).toBe(200)
    expect(response.data.requestSetUserEmail).toBe(true)
  }

  {
    const data = await rdb.GET(`send-SET_USER_EMAIL-${email}`)
    expect(data).not.toBeNull()
    const parsed = JSON.parse(data!)

    const response = await client.ApproveSetUserEmail({
      input: { token: parsed.token },
    })
    expect(response.status).toBe(200)
    expect(response.data.approveSetUserEmail.user.draftForm?.email).toBe(email)
  }

  {
    client.setToken(undefined)
    const response = await client.Viewer()
    expect(response.status).toBe(200)
    expect(response.data.viewer).toBeNull()
  }

  {
    const response = await client.Login({
      input: {
        password,
        username: email,
      },
    })
    expect(response.status).toBe(200)
    client.setToken(response.data.login.token)
  }

  {
    const response = await client.Viewer()
    expect(response.status).toBe(200)
    expect(response.data.viewer).not.toBeNull()
    expect(response.data.viewer).not.toBeUndefined()
    expect(response.data.viewer!.id).toBe(userId)
  }

  await rdb.quit()
})
