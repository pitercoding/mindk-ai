import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ChatPage from "@/pages/ChatPage";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <ChatPage />,
    },
]);

export default function AppRouter() {
    return <RouterProvider router={router} />;
}