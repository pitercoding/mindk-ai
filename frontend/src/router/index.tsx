import { createBrowserRouter, RouterProvider } from "react-router-dom";
import DashboardPage from "@/pages/DashboardPage";
import NotesPage from "@/pages/NotesPage";
import DashboardLayout from "@/layouts/DashboardLayout";

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
        ],
    },
]);

export default function AppRouter() {
    return <RouterProvider router={router} />;
}