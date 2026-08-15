import { createBrowserRouter, Navigate, RouterProvider } from "react-router-dom";

import DashboardPage from "@/pages/DashboardPage";
import NotesPage from "@/pages/NotesPage";

import DashboardLayout from "@/layouts/DashboardLayout";
import DocumentsPage from "@/pages/DocumentsPage";

export const router = createBrowserRouter([
    {
        path: "/",
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
]);

export default function AppRouter() {
    return <RouterProvider router={router} />;
}
