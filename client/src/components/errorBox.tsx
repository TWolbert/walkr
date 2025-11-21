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
        <div class="text-red-500 font-bold h-4 -mt-2">{s()!.message}</div>
      )}
    </>
  );
}
