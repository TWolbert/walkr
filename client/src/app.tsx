import { Suspense, type Component } from "solid-js";
import { A, useLocation } from "@solidjs/router";

const App: Component<{ children: Element }> = (props) => {
  const location = useLocation();

  return (
    <>
      <nav class="bg-gray-200 text-gray-900 px-4 fixed w-full">
        <ul class="flex justify-between">
          <div class="flex">
            <li class="py-2 px-4">
              <A href="/" class="no-underline hover:underline">
                Home
              </A>
            </li>
            <li class="py-2 px-4">
              <A href="/about" class="no-underline hover:underline">
                About
              </A>
            </li>
            <li class="py-2 px-4">
              <A href="/error" class="no-underline hover:underline">
                Error
              </A>
            </li>
          </div>
          <div class="flex">
            <li class="py-2 px-4">
              <A href="/auth/register" class="no-underline hover:underline">
                Register
              </A>
            </li>
            <li class="py-2 px-4">
              <A href="/auth/login" class="no-underline hover:underline">
                Login
              </A>
            </li>
          </div>
        </ul>
      </nav>

      <main class="h-screen">
        <Suspense>{props.children}</Suspense>
      </main>
    </>
  );
};

export default App;
