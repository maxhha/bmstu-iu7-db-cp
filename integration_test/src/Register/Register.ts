import { BaseTest, randid } from "../BaseTest"

export class Register extends BaseTest {
  protected password = `password-${randid()}`
  protected email = `test-email-${randid()}@email.com`
  protected userId!: string
  protected token!: string

  run() {
    it("should create new token on register", async () => {
      const response = await this.client.Register()
      expect(response.status).toBe(200)
      this.client.setToken(response.data.register.token)
    })

    it("should return non empty viewer for requests with authorization token", async () => {
      const response = await this.client.Viewer()
      expect(response.status).toBe(200)
      expect(response.data.viewer).not.toBeNull()
      expect(response.data.viewer).not.toBeUndefined()
      this.userId = response.data.viewer!.id
    })

    it("should set password", async () => {
      const response = await this.client.UpdateUserPassword({
        input: { password: this.password },
      })
      expect(response.status).toBe(200)
    })

    it("should request token to set new email", async () => {
      const response = await this.client.RequestSetUserEmail({
        input: { email: this.email },
      })
      expect(response.status).toBe(200)
      expect(response.data.requestSetUserEmail).toBe(true)
    })

    it("should set new email using token", async () => {
      const data = await this.rdb.GET(`send-SET_USER_EMAIL-${this.email}`)
      expect(data).not.toBeNull()
      const parsed = JSON.parse(data!)

      const response = await this.client.ApproveSetUserEmail({
        input: { token: parsed.token },
      })
      expect(response.status).toBe(200)
      expect(response.data.approveSetUserEmail.user.draftForm?.email).toBe(
        this.email
      )
    })
  }
}
