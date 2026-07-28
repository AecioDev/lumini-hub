// Paleta e helpers para avatares/tags "geradas" a partir de texto, já que o
// backend não guarda cor nenhuma para usuário/perfil (Role/User não têm
// campo de cor). Determinístico: o mesmo nome sempre cai na mesma cor.

const AVATAR_GRADIENTS = [
  "linear-gradient(135deg, #2563EB, #06B6D4)",
  "linear-gradient(135deg, #7C3AED, #2563EB)",
  "linear-gradient(135deg, #06B6D4, #7C3AED)",
  "linear-gradient(135deg, #1D4ED8, #7C3AED)",
  "linear-gradient(135deg, #06B6D4, #1D4ED8)",
];

const TAG_COLORS = ["blue", "purple", "cyan", "gold", "green", "magenta", "volcano"];

function hashString(value: string): number {
  let hash = 0;
  for (let i = 0; i < value.length; i++) {
    hash = (hash << 5) - hash + value.charCodeAt(i);
    hash |= 0;
  }
  return Math.abs(hash);
}

export function getInitials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "?";
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

export function getAvatarGradient(seed: string): string {
  return AVATAR_GRADIENTS[hashString(seed) % AVATAR_GRADIENTS.length];
}

export function getTagColor(seed: string): string {
  return TAG_COLORS[hashString(seed) % TAG_COLORS.length];
}
