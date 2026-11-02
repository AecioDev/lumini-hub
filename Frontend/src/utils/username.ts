// Palavras que ficam de fora ao escolher o último nome ("Maria da Silva" -> maria.silva).
const NAME_PARTICLES = new Set(["da", "das", "de", "di", "do", "dos", "du", "e"]);

// Sugere um nome de usuário a partir do nome completo: primeiro e último nome, sem acento, em minúsculas,
// separados por ponto ("Mônica Corvalan Machado" -> monica.machado). Com um nome só, usa ele. Devolve vazio
// quando o resultado seria curto demais para o usuário (mínimo 3 caracteres) ou passaria de 50.
export function suggestUsername(fullName: string): string {
  const words = fullName
    .normalize("NFD")
    .replace(/[̀-ͯ]/g, "")
    .toLowerCase()
    .split(/\s+/)
    .map((word) => word.replace(/[^a-z0-9]/g, ""))
    .filter(Boolean);
  if (words.length === 0) return "";

  const first = words[0];
  const rest = words.slice(1).filter((word) => !NAME_PARTICLES.has(word));
  const last = rest.length > 0 ? rest[rest.length - 1] : "";

  const username = last ? `${first}.${last}` : first;
  return username.length >= 3 && username.length <= 50 ? username : "";
}
