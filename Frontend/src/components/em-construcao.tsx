// Em: @/components/em-construcao.tsx (ou onde quer que ele esteja)

import { Construction } from "lucide-react";
import Link from "next/link";

export default function EmConstrucao() {
  return (
    <div className="flex flex-col items-center justify-center rounded-lg bg-gray-100 p-8 text-center shadow-sm">
      {" "}
      {/* Ajustes aqui */}
      <Construction className="mb-6 h-20 w-20 text-yellow-500" />
      <h1 className="mb-4 text-3xl font-bold text-gray-800">
        Área em construção {/* Título um pouco mais genérico para uma seção */}
      </h1>
      <p className="mb-6 text-lg text-gray-600">
        Estamos trabalhando para disponibilizar esta funcionalidade em breve.
        Agradecemos a sua paciência!
      </p>
      <Link href="/dashboard">
        <span className="inline-block rounded bg-blue-600 px-5 py-2 text-white transition hover:bg-blue-700">
          Voltar para o Dashboard
        </span>
      </Link>
    </div>
  );
}
