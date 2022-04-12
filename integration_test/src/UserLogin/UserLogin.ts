import { Register } from "../Register/Register"

export class UserLogin extends Register {
  login = () => {
    this.register()

    it("should return null viewer when authorization token not set", async () => {
      this.client.setToken(undefined)
      const response = await this.client.Viewer()
      expect(response.status).toBe(200)
      expect(response.data.viewer).toBeNull()
    })

    it("should return token on login", async () => {
      const response = await this.client.Login({
        input: {
          password: this.password,
          username: this.email,
        },
      })
      expect(response.status).toBe(200)
      this.client.setToken(response.data.login.token)
    })

    it("should return viewer with same id as was registred", async () => {
      const response = await this.client.Viewer()
      expect(response.status).toBe(200)
      expect(response.data.viewer).not.toBeNull()
      expect(response.data.viewer).not.toBeUndefined()
      expect(response.data.viewer!.id).toBe(this.userId)
    })
  }
}
