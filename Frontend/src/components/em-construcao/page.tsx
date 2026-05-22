import { Construction } from "lucide-react";
import Link from "next/link";

export default function EmConstrucao() {
  return (
    <div className="flex h-screen flex-col items-center justify-center bg-gray-100 px-4 text-center">
      <Construction className="h-20 w-20 text-yellow-500 mb-6" />
      <h1 className="text-4xl font-bold text-gray-800 mb-4">
        Página em construção
      </h1>
      <p className="text-lg text-gray-600 mb-6">
        Estamos trabalhando nisso! Em breve você poderá acessar essa
        funcionalidade.
      </p>
      <Link href="/dashboard">
        <span className="inline-block rounded bg-blue-600 px-5 py-2 text-white hover:bg-blue-700 transition">
          Voltar para o Dashboard
        </span>
      </Link>
    </div>
  );
}
