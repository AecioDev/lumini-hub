import { z } from "zod";

export const legacyLocationSchema = z.object({
  cod_loc_am: z.number().int(),
  des_loc_am: z.string(),
  v_mais_local_empresa: z.number().int(),
});

export type LegacyLocation = z.infer<typeof legacyLocationSchema>;

export const legacyCompanySchema = z.object({
  cod_cencus: z.number().int(),
  des_cencus: z.string(),
  cgc_cencus: z.string(),
});

export type LegacyCompany = z.infer<typeof legacyCompanySchema>;
