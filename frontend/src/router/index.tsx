import { createBrowserRouter, RouterProvider } from "react-router-dom";

import DashboardPage from "@/pages/DashboardPage";
import NotesPage from "@/pages/NotesPage";
import AuthPage from "@/pages/AuthPage";
import NotFoundPage from "@/pages/NotFoundPage";

import DashboardLayout from "@/layouts/DashboardLayout";
import DocumentsPage from "@/pages/DocumentsPage";

import ProtectedRoute from "@/components/auth/ProtectedRoute";

const router = createBrowserRouter([
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
                        element: <NotFoundPage />,
                    },
                ],
            },
        ],
    },
]);

export default function AppRouter() {
    return <RouterProvider router={router} />;
}
