import { redirect } from "next/navigation";

export default function Home() {
  // Redirect to dashboard if authenticated, otherwise to login
  redirect("/login");
  return null;
}
