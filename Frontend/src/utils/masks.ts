// Helpers de máscara pra CPF/CNPJ/telefone — mascara na tela, mas só os
// dígitos brutos vão pro backend (ver "Masked fields" em CLAUDE.md). Sem
// biblioteca de máscara instalada ainda no projeto (checado 2026-09-05),
// então é um helper local simples em vez de regex por lib.

export function onlyDigits(value: string): string {
  return value.replace(/\D/g, "");
}

export function maskCPF(value: string): string {
  const digits = onlyDigits(value).slice(0, 11);
  return digits
    .replace(/(\d{3})(\d)/, "$1.$2")
    .replace(/(\d{3})(\d)/, "$1.$2")
    .replace(/(\d{3})(\d{1,2})$/, "$1-$2");
}

export function maskCNPJ(value: string): string {
  const digits = onlyDigits(value).slice(0, 14);
  return digits
    .replace(/(\d{2})(\d)/, "$1.$2")
    .replace(/(\d{3})(\d)/, "$1.$2")
    .replace(/(\d{3})(\d)/, "$1/$2")
    .replace(/(\d{4})(\d{1,2})$/, "$1-$2");
}

// Formata um documento (CPF ou CNPJ) pra exibição, decidindo o formato pela
// quantidade de dígitos — usado pro documento extraído do certificado
// digital, que pode ser pessoa física ou jurídica.
export function maskDocument(value: string): string {
  const digits = onlyDigits(value);
  if (digits.length > 11) return maskCNPJ(digits);
  if (digits.length > 0) return maskCPF(digits);
  return value;
}

// Telefone fixo (10 dígitos) ou celular (11 dígitos) — decide o formato
// dinamicamente conforme o usuário digita.
export function maskPhone(value: string): string {
  const digits = onlyDigits(value).slice(0, 11);
  if (digits.length <= 10) {
    return digits
      .replace(/(\d{2})(\d)/, "($1) $2")
      .replace(/(\d{4})(\d{1,4})$/, "$1-$2");
  }
  return digits
    .replace(/(\d{2})(\d)/, "($1) $2")
    .replace(/(\d{5})(\d{1,4})$/, "$1-$2");
}
