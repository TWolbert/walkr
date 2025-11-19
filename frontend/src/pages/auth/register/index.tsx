import { createSignal } from "solid-js";
import { DOMElement } from "solid-js/jsx-runtime";
import { api, User } from "../../../api/api";

type Status =
  | {
      type: "error" | "success";
      message: string;
      for: string;
    }
  | undefined;

export default function Register() {
  const [username, setUsername] = createSignal("");
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [status, setStatus] = createSignal<Status>();
  const [submitting, setSubmitting] = createSignal(false);

  async function submitForm(
    e: SubmitEvent & { currentTarget: HTMLFormElement; target: DOMElement },
  ) {
    e.preventDefault();
    setSubmitting(true);
    setStatus(undefined);

    try {
      const res = await api.post("user", {
        body: JSON.stringify({
          username: username(),
          email: email(),
          password: password(),
        }),
        headers: {
          "content-type": "application/json",
        },
      });

      const data = (await res.json()) as
        | User
        | { error: string; detail: string; for: string };

      setSubmitting(false);

      if ("error" in data) {
        setStatus({
          type: "error",
          message: data.detail,
          for: data.for,
        });
        return;
      }

      setStatus({
        type: "success",
        message: "Account created successfully",
        for: "any",
      });
    } catch (err) {
      setSubmitting(false);
      setStatus({
        type: "error",
        message: "Something went wrong while creating an account!",
        for: "any",
      });
    }
  }

  return (
    <section class="p-3 flex items-center justify-center h-full">
      <form
        class="flex flex-col gap-2 border p-3 min-w-[30vw] rounded-xl shadow"
        method="post"
        onSubmit={submitForm}
      >
        <label for="username">Username</label>
        <input
          type="text"
          name="username"
          value={username()}
          id="username"
          placeholder="Your username..."
          class="px-3 py-2 rounded shadow-sm"
          onInput={(e) => setUsername(e.currentTarget.value)}
        />
        <ErrorBox status={status} forKey="username" />

        <label for="email">Email</label>
        <input
          type="email"
          name="email"
          value={email()}
          id="email"
          placeholder="Your email..."
          class="px-3 py-2 rounded shadow-sm"
          onInput={(e) => setEmail(e.currentTarget.value)}
        />
        <ErrorBox status={status} forKey="email" />

        <label for="password">Password</label>
        <input
          type="password"
          name="password"
          id="password"
          value={password()}
          placeholder="Your password..."
          class="px-3 py-2 rounded shadow-sm"
          onInput={(e) => setPassword(e.currentTarget.value)}
        />
        <ErrorBox status={status} forKey="password" />

        <button
          class="px-3 py-2 bg-orange-500 text-white font-bold rounded-md shadow-md disabled:bg-orange-300 transition-all"
          type="submit"
          disabled={submitting()}
        >
          Submit
        </button>

        <ErrorBox status={status} forKey="any" />
        {status()?.type === "success" && (
          <span class="font-bold text-green-500">{status()?.message}</span>
        )}
      </form>
    </section>
  );
}

// child
export function ErrorBox({
  status,
  forKey,
}: {
  status?: () =>
    | { type: "error" | "success"; message: string; for: string }
    | undefined;
  forKey: string;
}) {
  const s = () => status?.(); // read accessor reactively
  return (
    <>
      {s() && s()!.type === "error" && s()!.for === forKey && (
        <div class="text-red-500 font-bold">{s()!.message}</div>
      )}
    </>
  );
}
