import { createBrowserRouter, RouterProvider } from "react-router-dom";

import DashboardPage from "@/pages/DashboardPage";
import NotesPage from "@/pages/NotesPage";
import ChatPage from "@/pages/ChatPage";

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
            {
                path: "chat",
                element: <ChatPage />,
            },
        ],
    },
]);

export default function AppRouter() {
    return <RouterProvider router={router} />;
}
