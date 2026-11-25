import { useEffect, useState } from "react";
import { Outlet, useLocation } from "react-router-dom";

import Sidebar from "@/components/layout/Sidebar";
import Header from "@/components/layout/Header";

export default function DashboardLayout() {

    const { pathname } = useLocation();

    // The drawer is tied to the route it was opened on, so navigating
    // to another page closes it without any effect.
    const [openedAtPath, setOpenedAtPath] = useState<string | null>(null);

    const isSidebarOpen = openedAtPath === pathname;

    function closeSidebar() {
        setOpenedAtPath(null);
    }

    function toggleSidebar() {
        setOpenedAtPath(isSidebarOpen ? null : pathname);
    }

    // Close the drawer with the Escape key.
    useEffect(() => {

        if (!isSidebarOpen) {
            return;
        }

        function handleKeyDown(event: KeyboardEvent) {
            if (event.key === "Escape") {
                closeSidebar();
            }
        }

        window.addEventListener("keydown", handleKeyDown);

        return () => {
            window.removeEventListener("keydown", handleKeyDown);
        };

    }, [isSidebarOpen]);

    return (
        <div className="dashboard-layout">

            <Sidebar
                isOpen={isSidebarOpen}
                onClose={closeSidebar}
            />

            {isSidebarOpen && (
                <div
                    className="sidebar-backdrop"
                    aria-hidden="true"
                    onClick={closeSidebar}
                />
            )}

            <div className="dashboard-content">

                <Header
                    isMenuOpen={isSidebarOpen}
                    onMenuClick={toggleSidebar}
                />

                <main>
                    <Outlet />
                </main>

            </div>

        </div>
    );
}
