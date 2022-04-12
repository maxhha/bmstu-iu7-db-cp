import { BaseTest, randid } from "../BaseTest"
import { UserFormState } from "../query"

export class UserFormModeration extends BaseTest {
  private password = "password" + randid()
  private email = `test-email-${randid()}@email.com`
  private name = "test-client" + randid()
  private phone = "test-phone" + randid()
  private userId!: string
  private userToken!: string

  async before() {
    await super.before()

    const { userId, token: userToken } = await this.register()
    this.userId = userId
    this.userToken = userToken

    await this.fillUser({
      token: userToken,
      email: this.email,
      phone: this.phone,
      password: this.password,
      form: { name: this.name },
    })
  }

  run = () => {
    it("should request for draft form moderation", async () => {
      await this.client.RequestModerateUserForm().then((response) => {
        expect(response.status).toBe(200)
        expect(response.data.requestModerateUserForm).toBe(true)
      })
    })

    it("should approve with token draft form moderation", async () => {
      const data = await this.rdb.GET(`send-MODERATE_USER_FORM-${this.phone}`)
      expect(data).not.toBeNull()
      const parsed = JSON.parse(data!)

      await this.client
        .ApproveModerateUserForm({ input: { token: parsed.token } })
        .then((response) => {
          expect(response.status).toBe(200)
          expect(
            response.data.approveModerateUserForm.user.draftForm?.state
          ).toBe(UserFormState.Moderating)
        })
    })
  }
}
