import { createBrowserRouter, Navigate, RouterProvider } from "react-router-dom";

import DashboardPage from "@/pages/DashboardPage";
import NotesPage from "@/pages/NotesPage";
import AuthPage from "@/pages/AuthPage";

import DashboardLayout from "@/layouts/DashboardLayout";
import DocumentsPage from "@/pages/DocumentsPage";

import ProtectedRoute from "@/components/auth/ProtectedRoute";

export const router = createBrowserRouter([
    {
        path: "/auth",
        element: <AuthPage />,
    },
    {
        path: "/",
        element: <ProtectedRoute />,
        children: [
            {
                element: <DashboardLayout />,
                children: [
                    {
                        index: true,
                        element: <DashboardPage />,
                    },
                    {
                        path: "notes",
                        element: <NotesPage />,
                    },
                    {
                        path: "documents",
                        element: <DocumentsPage />,
                    },
                    {
                        path: "*",
                        element: <Navigate to="/" replace />,
                    },
                ],
            },
        ],
    },
]);

export default function AppRouter() {
    return <RouterProvider router={router} />;
}
