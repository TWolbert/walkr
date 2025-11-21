import { Setter, Show } from "solid-js";
import { ErrorBox } from "./errorBox";

export function ValidatedInput({
  value,
  name,
  placeHolder,
  setValue,
  status,
  label,
  type,
  required = false,
}: {
  value?: () => string | undefined;
  name: string;
  placeHolder: string;
  setValue: Setter<string>;
  status:
    | (() =>
        | {
            type: "error" | "success";
            message: string;
            for: string;
          }
        | undefined)
    | undefined;
  label: string;
  type: string;
  required?: boolean | null;
}) {
  const val = () => value?.();
  const s = () => status?.();
  return (
    <>
      <label for={name} class="flex gap-1 items-center">
        <Show when={required}>
          <Show
            when={s() && s()!.type === "error" && s()!.for === name}
            fallback={
              <p class=" text-red-500 font-bold w-4 h-4 flex items-center justify-center rounded-full">
                *
              </p>
            }
          >
            <p class="bg-red-500 text-white font-bold w-4 h-4 flex items-center justify-center rounded-full">
              !
            </p>
          </Show>
        </Show>
        {label}
      </label>
      <input
        name={name}
        value={val()}
        id={name}
        placeholder={placeHolder}
        class={
          "px-3 py-2 rounded border border-transparent shadow-sm" +
          (s() && s()!.type === "error" && s()!.for === name
            ? "  border-red-500!"
            : "")
        }
        onInput={(e) => setValue(e.currentTarget.value)}
        type={type || "text"}
      />
      <Show
        when={s() && s()!.type === "error" && s()!.for === name}
        fallback={<p class="font-bold text-zinc-400 h-4 -mt-2"> </p>}
      >
        <ErrorBox status={s} forKey={name} />
      </Show>
    </>
  );
}
