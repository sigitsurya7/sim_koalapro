export function formatTanggalID(iso: string, withJam = false) {
  const d = new Date(iso);

  const tanggal = new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  }).format(d);

  if (!withJam) {
    return tanggal;
  }

  const jam = new Intl.DateTimeFormat("id-ID", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
    // kalau ingin jam WIB konsisten, pakai timeZone: "Asia/Jakarta"
    timeZone: "Asia/Jakarta",
  }).format(d);

  return `${tanggal} ${jam}`;
}