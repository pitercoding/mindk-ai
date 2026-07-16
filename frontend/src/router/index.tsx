import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ChatPage from "@/pages/ChatPage";
import NotesPage from "@/pages/NotesPage";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <ChatPage />,
    },
    {
        path: "/notes",
        element: <NotesPage />,
    }
]);

export default function AppRouter() {
    return <RouterProvider router={router} />;
}