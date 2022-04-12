import { UserFormModeration } from "./UserFormModeration"

describe("UserFormModeration", () => {
  const test = new UserFormModeration()
  beforeAll(() => test.before())
  afterAll(() => test.after())

  test.run()
})
