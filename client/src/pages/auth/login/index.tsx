import { createSignal, Show } from "solid-js";
import { DOMElement } from "solid-js/jsx-runtime";
import { api, User } from "../../../api/api";
import { ErrorBox } from "../../../components/errorBox";
import { ValidatedInput } from "../../../components/validatedInput";

type Status =
  | {
      type: "error" | "success";
      message: string;
      for: string;
    }
  | undefined;

export default function Register() {
  const [usernameOrEmail, setUsernameOrEmail] = createSignal("");
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
      const res = await api.post("user/login", {
        body: JSON.stringify({
          username_or_email: usernameOrEmail(),
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
        message: "Logged in succesfully",
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
        class="min-w-full md:min-w-[50vw] flex flex-col gap-2 border p-3 rounded-xl shadow"
        method="post"
        onSubmit={submitForm}
      >
        <ValidatedInput
          value={usernameOrEmail}
          setValue={setUsernameOrEmail}
          name="username_or_email"
          status={status}
          label="Username or Email"
          placeHolder="Your username or email..."
          type="text"
          required
        />

        <ValidatedInput
          value={password}
          setValue={setPassword}
          name="password"
          status={status}
          label="Password"
          placeHolder="Your password..."
          type="password"
          required
        />

        <button
          class="px-3 py-2 bg-orange-500 text-white font-bold rounded-md shadow-md disabled:bg-orange-300 transition-all"
          type="submit"
          disabled={submitting()}
        >
          <Show when={submitting()} fallback={<p>Log in</p>}>
            <p>Logging in...</p>
          </Show>
        </button>

        <ErrorBox status={status} forKey="any" />
        {status()?.type === "success" && (
          <span class="font-bold text-green-500">{status()?.message}</span>
        )}
      </form>
    </section>
  );
}
