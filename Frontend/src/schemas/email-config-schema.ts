import { z } from "zod";

// Máximos iguais ao size das colunas (SMTPSettings, no api.core). Ligado, servidor e remetente são
// obrigatórios; desligado, dá para salvar pela metade.
export const emailConfigSchema = z
  .object({
    host: z.string().max(150, "Máximo 150 caracteres"),
    port: z
      .number({ message: "Informe a porta" })
      .int("Porta inválida")
      .min(1, "Porta inválida")
      .max(65535, "Porta inválida"),
    security: z.enum(["starttls", "ssl", "none"], { message: "Selecione a segurança" }),
    username: z.string().max(150, "Máximo 150 caracteres"),
    password: z.string().max(255, "Máximo 255 caracteres"),
    from_name: z.string().max(100, "Máximo 100 caracteres"),
    from_email: z
      .string()
      .max(150, "Máximo 150 caracteres")
      .refine((value) => value === "" || z.string().email().safeParse(value).success, "E-mail inválido"),
    enabled: z.boolean(),
  })
  .superRefine((values, ctx) => {
    if (!values.enabled) return;
    if (values.host.trim() === "") {
      ctx.addIssue({ code: "custom", path: ["host"], message: "Informe o servidor para ativar o envio" });
    }
    if (values.from_email.trim() === "") {
      ctx.addIssue({ code: "custom", path: ["from_email"], message: "Informe o e-mail do remetente" });
    }
  });

export type EmailConfigFormValues = z.infer<typeof emailConfigSchema>;

export const testEmailSchema = z.string().email("E-mail inválido").max(150, "Máximo 150 caracteres");
