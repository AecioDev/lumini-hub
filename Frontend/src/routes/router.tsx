import { createBrowserRouter, Navigate } from "react-router-dom";
import { AuthLayout } from "./AuthLayout";
import { AppLayout } from "./AppLayout";
import { ProtectedRoute } from "./ProtectedRoute";
import { RequirePermission } from "./RequirePermission";
import { PlaceholderPage } from "@/pages/PlaceholderPage";
import { LoginPage } from "@/pages/auth/LoginPage";
import { ForgotPasswordPage } from "@/pages/auth/ForgotPasswordPage";
import { CrmDashboardPage } from "@/pages/dashboard/CrmDashboardPage";
import { UsersListPage } from "@/pages/users/UsersListPage";
import { CreateUserPage } from "@/pages/users/CreateUserPage";
import { EditUserPage } from "@/pages/users/EditUserPage";
import { IntegrationsPage } from "@/pages/settings/IntegrationsPage";
import { RolesAndPermissionsPage } from "@/pages/settings/RolesAndPermissionsPage";
import { CreateRolePage } from "@/pages/settings/roles/CreateRolePage";
import { EditRolePage } from "@/pages/settings/roles/EditRolePage";
import { CreatePermissionPage } from "@/pages/settings/permissions/CreatePermissionPage";
import { EditPermissionPage } from "@/pages/settings/permissions/EditPermissionPage";
import { CompaniesListPage } from "@/pages/settings/companies/CompaniesListPage";
import { CreateCompanyPage } from "@/pages/settings/companies/CreateCompanyPage";
import { EditCompanyPage } from "@/pages/settings/companies/EditCompanyPage";

export const router = createBrowserRouter([
  {
    element: <AuthLayout />,
    children: [
      { path: "/login", element: <LoginPage /> },
      { path: "/forgot-password", element: <ForgotPasswordPage /> },
    ],
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <AppLayout />,
        children: [
          { index: true, element: <Navigate to="/dashboard" replace /> },
          { path: "dashboard", element: <CrmDashboardPage /> },
          {
            // Precisa bater com o href real semeado em menu_items ("Empresas" -> /settings/companies).
            path: "settings/companies",
            children: [
              {
                element: <RequirePermission permission="companies.view" />,
                children: [{ index: true, element: <CompaniesListPage /> }],
              },
              {
                element: <RequirePermission permission="companies.create" />,
                children: [{ path: "create", element: <CreateCompanyPage /> }],
              },
              {
                element: <RequirePermission permission="companies.edit" />,
                children: [{ path: ":id/edit", element: <EditCompanyPage /> }],
              },
            ],
          },
          {
            // Precisa bater com o href real semeado em menu_items ("Usuários" -> /settings/users).
            path: "settings/users",
            children: [
              {
                element: <RequirePermission permission="users.view" />,
                children: [{ index: true, element: <UsersListPage /> }],
              },
              {
                element: <RequirePermission permission="users.create" />,
                children: [{ path: "create", element: <CreateUserPage /> }],
              },
              {
                element: <RequirePermission permission="users.edit" />,
                children: [{ path: ":id/edit", element: <EditUserPage /> }],
              },
            ],
          },
          {
            // Precisa bater com o href real semeado em menu_items ("Integrações" -> /settings/integrations).
            path: "settings/integrations",
            element: <RequirePermission permission="integrations.view" />,
            children: [{ index: true, element: <IntegrationsPage /> }],
          },
          {
            // Precisa bater com o href real semeado em menu_items
            // ("Perfis e Permissões" -> /settings/roles, gate roles.view — mesma
            // permissão que já abre o menu pro ADMIN; admin.create_permissions
            // fica reservado pro CRUD de /menu-items, dev-only).
            path: "settings/roles",
            children: [
              {
                element: <RequirePermission permission="roles.view" />,
                children: [{ index: true, element: <RolesAndPermissionsPage /> }],
              },
              {
                element: <RequirePermission permission="roles.create" />,
                children: [{ path: "create", element: <CreateRolePage /> }],
              },
              {
                element: <RequirePermission permission="roles.edit" />,
                children: [{ path: ":id/edit", element: <EditRolePage /> }],
              },
            ],
          },
          {
            path: "settings/permissions",
            children: [
              {
                element: <RequirePermission permission="permissions.create" />,
                children: [{ path: "create", element: <CreatePermissionPage /> }],
              },
              {
                element: <RequirePermission permission="permissions.edit" />,
                children: [{ path: ":id/edit", element: <EditPermissionPage /> }],
              },
            ],
          },
          { path: "*", element: <PlaceholderPage /> },
        ],
      },
    ],
  },
], {
  // Opt-in antecipado nas flags da v7 pra sumir com os warnings de depreciação
  // no console (react-router-dom 6.28 já implementa o comportamento, só
  // precisa habilitar). v7_startTransition é flag de <RouterProvider>, não
  // do router em si — vai em App.tsx.
  future: {
    v7_relativeSplatPath: true,
    v7_fetcherPersist: true,
    v7_normalizeFormMethod: true,
    v7_partialHydration: true,
    v7_skipActionErrorRevalidation: true,
  },
});
