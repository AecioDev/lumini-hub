// Paleta e helpers para avatares/tags "geradas" a partir de texto, já que o
// backend não guarda cor nenhuma para usuário/perfil (Role/User não têm
// campo de cor). Determinístico: o mesmo nome sempre cai na mesma cor.

const AVATAR_GRADIENTS = [
  "linear-gradient(135deg, #1677ff, #9254de)",
  "linear-gradient(135deg, #13a8a8, #1677ff)",
  "linear-gradient(135deg, #eb2f96, #722ed1)",
  "linear-gradient(135deg, #fa8c16, #d4b106)",
  "linear-gradient(135deg, #52c41a, #13a8a8)",
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
