import { UserLogin } from "./UserLogin"

describe("UserLogin", () => {
  const test = new UserLogin()
  beforeAll(test.before)
  afterAll(test.after)

  test.login()
})

// describe("UserLogin", () => {
//   const password = "HelloWorld-password"
//   const email = "test-email@email.com"
//   const client = getClient()
//   let rdb: Awaited<ReturnType<typeof getRedisClient>>

//   beforeAll(async () => {
//     rdb = await getRedisClient()
//   })

//   afterAll(() => rdb.quit())

//   it("should create new token on register", async () => {
//     const response = await client.Register()
//     expect(response.status).toBe(200)
//     client.setToken(response.data.register.token)
//   })

//   let userId: string
//   it("should return non empty viewer for requests with authorization token", async () => {
//     const response = await client.Viewer()
//     expect(response.status).toBe(200)
//     expect(response.data.viewer).not.toBeNull()
//     expect(response.data.viewer).not.toBeUndefined()
//     userId = response.data.viewer!.id
//   })

//   it("should set password", async () => {
//     const response = await client.UpdateUserPassword({ input: { password } })
//     expect(response.status).toBe(200)
//   })

//   it("should request token to set new email", async () => {
//     const response = await client.RequestSetUserEmail({ input: { email } })
//     expect(response.status).toBe(200)
//     expect(response.data.requestSetUserEmail).toBe(true)
//   })

//   it("should set new email using token", async () => {
//     const data = await rdb.GET(`send-SET_USER_EMAIL-${email}`)
//     expect(data).not.toBeNull()
//     const parsed = JSON.parse(data!)

//     const response = await client.ApproveSetUserEmail({
//       input: { token: parsed.token },
//     })
//     expect(response.status).toBe(200)
//     expect(response.data.approveSetUserEmail.user.draftForm?.email).toBe(email)
//   })

//   it("should return null viewer when authorization token not set", async () => {
//     client.setToken(undefined)
//     const response = await client.Viewer()
//     expect(response.status).toBe(200)
//     expect(response.data.viewer).toBeNull()
//   })

//   it("should return token on login", async () => {
//     const response = await client.Login({
//       input: {
//         password,
//         username: email,
//       },
//     })
//     expect(response.status).toBe(200)
//     client.setToken(response.data.login.token)
//   })

//   it("should return viewer with same id as was registred", async () => {
//     const response = await client.Viewer()
//     expect(response.status).toBe(200)
//     expect(response.data.viewer).not.toBeNull()
//     expect(response.data.viewer).not.toBeUndefined()
//     expect(response.data.viewer!.id).toBe(userId)
//   })
// })
