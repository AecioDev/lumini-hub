import { Metadata } from "next";
import { OpenGraph } from "next/dist/lib/metadata/types/opengraph-types";

enum MODE {
  DARK = "dark",
  LIGHT = "light",
}

export const siteConfig = {
  title: "Simple ERP - Gestão Empresarial",
  description: `Simple ERP - Gestão Empresarial é um sistema de gestão integrado, desenvolvido para simplificar os processos operacionais de micro e pequenas empresas. A plataforma oferece recursos essenciais como controle de clientes, finanças, vendas e estoque, com uma interface moderna, responsiva e fácil de usar. Ideal para empresas que buscam centralizar suas informações e ganhar eficiência na administração do negócio.`,
  logo: "/logo.svg",
  icon: "/logo-short.svg",
  mode: MODE.LIGHT,
  // TODO: favicon
};

export const metaObject = (
  title?: string,
  openGraph?: OpenGraph,
  description: string = siteConfig.description
): Metadata => {
  return {
    title: title
      ? `${title} - Simple ERP - Gestão Empresarial`
      : siteConfig.title,
    description,
    openGraph: openGraph ?? {
      title: title ? `${title} - Simple ERP - Gestão Empresarial` : title,
      description,
      url: "https://simpleerp.com.br",
      siteName: "Simple ERP - Gestão Empresarial", // https://developers.google.com/search/docs/appearance/site-names
      images: {
        url: "https://simpleerp.com.br/og-image.png",
        width: 1200,
        height: 630,
      },
      locale: "pt_BR",
      type: "website",
    },
  };
};
