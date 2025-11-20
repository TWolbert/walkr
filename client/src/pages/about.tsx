import { createEffect, createSignal, Suspense } from "solid-js";
import AboutData from "./about.data";
import type { User } from "../api/api";

export default function About() {
  const data = AboutData();

  const [error, setError] = createSignal("");
  const [user, setUser] = createSignal<User>();

  createEffect(() => {
    const v = data();
    if (!v) return;
    if (!("ID" in v)) {
      setError(v.error ?? "Unknown error");
    } else {
      setUser(v);
    }
  });

  return (
    <section class="bg-pink-100 text-gray-700 p-8">
      <h1 class="text-2xl font-bold">About</h1>

      <p class="mt-4">A page all about this website.</p>

      <p>
        <span>We love</span>
        <Suspense fallback={<span>...</span>}>
          <span>
            &nbsp;
            {error()}
          </span>
          <span>&nbsp; {JSON.stringify(user()!)}</span>
        </Suspense>
      </p>
    </section>
  );
}
