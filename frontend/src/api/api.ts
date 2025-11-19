import ky from "ky";

export const api = ky.create({
  prefixUrl: "/api",
  throwHttpErrors: false,
});

export async function test() {
  return await (await api.get("users")).json();
}

export interface User {
  ID: number;
  Username: string;
  Email: string;
  Password: string;
  CreatedAt: string;
  UpdatedAt: string;
}

export async function createTestUser() {
  try {
    const json = await api
      .post("user", {
        json: {
          username: "testuser234",
          email: "testuser234@goed.email",
          password: "testtest",
        },
      })
      .json<User>();
    return json;
  } catch (e: any) {
    const err = e?.response ? await e.response.json().catch(() => null) : null;
    if (err && typeof err.error === "string") return { error: err.error };
    return { error: "Request failed" };
  }
}
