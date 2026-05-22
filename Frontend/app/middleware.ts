import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PROTECTED_ROUTES = ["/dashboard", "/cadastros", "/estoque", "/usuarios"];
const PUBLIC_AUTH_ROUTES = ["/login", "/signin"];

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const accessToken = request.cookies.get("access_token")?.value;

  if (accessToken && PUBLIC_AUTH_ROUTES.includes(pathname)) {
    return NextResponse.redirect(new URL("/dashboard", request.url));
  }

  const isProtectedRoute = PROTECTED_ROUTES.some((route) =>
    pathname.startsWith(route)
  );

  if (isProtectedRoute) {
    if (!accessToken) {
      const loginUrl = new URL("/login", request.url);
      loginUrl.searchParams.set("redirect", pathname);
      return NextResponse.redirect(loginUrl);
    }

    try {
      // Decodifica o JWT para verificar a expiração (NÃO valida a assinatura aqui)
      const tokenPayload = parseJwt(accessToken);
      if (!tokenPayload || tokenPayload.exp * 1000 < Date.now()) {
        console.log("Access Token expirado ou inválido no middleware");
        const response = NextResponse.redirect(new URL("/login", request.url));
        response.cookies.delete("access_token");
        response.cookies.delete("refresh_token");
        return response;
      }
    } catch (error) {
      console.error("Erro ao processar token no middleware:", error);
      const response = NextResponse.redirect(new URL("/login", request.url));
      response.cookies.delete("access_token");
      response.cookies.delete("refresh_token");
      return response;
    }
  }
  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};

function parseJwt(token: string) {
  try {
    return JSON.parse(atob(token.split(".")[1]));
  } catch (e) {
    return null;
  }
}
